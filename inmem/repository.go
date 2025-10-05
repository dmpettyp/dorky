package inmem

import (
	"errors"
	"slices"

	"github.com/dmpettyp/dorky"
)

// entity defines the generic interface for entities that can be stored in an
// in-memory repository. All entities that can be stored in an in-memory
// repository must implement a:
// - GetEvents() method that returns event messages generated with the entity
// - Clone() method that returns a copy of the entity
type entity[T any] interface {
	GetEvents() []dorky.Event
	ResetEvents()
	Clone() T
}

// Repository defines a generic in-memory repository that can be used as a base
// for in-memory repositories for concrete entities.
//
// Entities stored in the repository must implement the generic entity
// interface
type Repository[Entity entity[Entity]] struct {
	// Entities is the collection of "persisted". This is the durable backing for
	// the inmen repository.
	Entities []Entity

	// Transaction contains all entities accessed or added during this UoW transaction.
	// These entities may be modified and will be persisted on Save() or discarded on Reset()
	Transaction []Entity

	identityEqualFn   func(Entity, Entity) bool
	constraintEqualFn func(Entity, Entity) bool
}

func CreateRepository[Entity entity[Entity]](
	identityEqualFn func(Entity, Entity) bool,
	constraintEqualFn func(Entity, Entity) bool,
) (
	Repository[Entity],
	error,
) {
	if identityEqualFn == nil {
		return Repository[Entity]{}, errors.New("identityEqualFn cannot be nil")
	}
	if constraintEqualFn == nil {
		return Repository[Entity]{}, errors.New("constraintEqualFn cannot be nil")
	}

	return Repository[Entity]{
		Entities:          nil,
		Transaction:       nil,
		identityEqualFn:   identityEqualFn,
		constraintEqualFn: constraintEqualFn,
	}, nil
}

// Add verifies that an equivalent entity doesn't already exist in the
// repository and then adds it to the repository's uncommitted entities.
//
// An error with code CodeAlreadyExists will be returned if the entity
// duplicates one that's already persisted or has uncommitted changes.
func (repo *Repository[Entity]) Add(
	toAdd Entity,
) error {
	for _, entity := range append(repo.Entities, repo.Transaction...) {
		if repo.constraintEqualFn(toAdd, entity) {
			return ErrAlreadyExists
		}
	}

	repo.Transaction = append(repo.Transaction, toAdd)

	return nil
}

// FindOne uses the provided match function to look for an entity in the
// repository to be returned.
//
// The entity returned is added to the repository's Transaction collection.
// The returned entity can be directly modified, but those changes won't
// persist to Entities until Save is called.
func (repo *Repository[Entity]) FindOne(
	matchFn func(Entity) bool,
) (
	Entity,
	error,
) {
	for _, entity := range repo.Transaction {
		if matchFn(entity) {
			return entity, nil
		}
	}

	for _, entity := range repo.Entities {
		if !matchFn(entity) {
			continue
		}

		// Transaction entity is a copy so that writes don't affect the persisted Entities
		// until the repo is saved
		transactionEntity := entity.Clone()

		repo.Transaction = append(repo.Transaction, transactionEntity)

		return transactionEntity, nil
	}

	var zero Entity
	return zero, ErrNotFound
}

// FindAll uses the provided match function to look for all entities in the
// repository to be returned.
//
// All entities returned are added to the repository's Transaction collection.
// The returned entities can be directly modified, but those changes won't
// persist to Entities until Save is called.
func (repo *Repository[Entity]) FindAll(
	matchFn func(Entity) bool,
) (
	[]Entity,
	error,
) {
	var all []Entity

	for _, entity := range repo.Transaction {
		if matchFn(entity) {
			all = append(all, entity)
		}
	}

	for _, entity := range repo.Entities {
		if !matchFn(entity) {
			continue
		}

		// Check if an entity with matching identity is already in all
		if slices.ContainsFunc(all, func(e Entity) bool {
			return repo.identityEqualFn(e, entity)
		}) {
			continue
		}

		transactionEntity := entity.Clone()
		repo.Transaction = append(repo.Transaction, transactionEntity)

		all = append(all, transactionEntity)
	}

	return all, nil
}

// Save persists any entities found in the repository's Transaction collection into
// its persistent store (the Entities collection)
func (repo *Repository[Entity]) Save() ([]dorky.Event, error) {
	// Validate all constraints before persisting anything to ensure atomicity
	for _, toSave := range repo.Transaction {
		for _, entity := range repo.Entities {
			if repo.identityEqualFn(entity, toSave) {
				continue
			}

			if repo.constraintEqualFn(entity, toSave) {
				return nil, ErrAlreadyExists
			}
		}
	}

	persist := func(toSave Entity) {
		// Replace the entity whos identity matches
		for idx, entity := range repo.Entities {
			if repo.identityEqualFn(entity, toSave) {
				repo.Entities[idx] = toSave
				return
			}
		}

		// No identiy match found, add the entity as a new entity
		repo.Entities = append(repo.Entities, toSave)
	}

	var events []dorky.Event

	for _, transactionEntity := range repo.Transaction {
		persist(transactionEntity)
		events = append(events, transactionEntity.GetEvents()...)
		transactionEntity.ResetEvents()
	}

	repo.Transaction = nil

	return events, nil
}

// Reset is used to clear the repository's Transaction collection and any changes made
// to the transaction entities are forgotten and can't be saved.
func (repo *Repository[Entity]) Reset() {
	repo.Transaction = nil
}
