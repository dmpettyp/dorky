package inmem

import (
	"context"

	"github.com/dmpettyp/dorky"
)

type repo interface {
	Save() ([]dorky.Event, error)
	Reset()
}

// InmemUnitOfWork provides a generic implementation of an in-memory UnitOfWork
// with arbitrary in-memory Repositories that can be embedded within a specific
// in-memory UnitOfWork with concrete in-memory Repositories
type UnitOfWork[Repos any] struct {
	repos    Repos
	repoList []repo
}

func NewUnitOfWork[Repos any](repos Repos, repoList ...repo) *UnitOfWork[Repos] {
	unitOfWork := &UnitOfWork[Repos]{
		repos:    repos,
		repoList: repoList,
	}
	return unitOfWork
}

// Run executes f in a transaction and returns the events that were created by executing f. If f returns an
// error then the transaction will be rolled back.
func (uow *UnitOfWork[Repos]) Run(
	_ context.Context,
	f func(Repos) error,
) ([]dorky.Event, error) {
	resetRepos := func() {
		for _, r := range uow.repoList {
			r.Reset()
		}
	}

	resetRepos()

	// always reset the in-mem working set of entities after the unit of work
	// completes (whether success or failure) to prepare for the next transaction
	defer resetRepos()

	err := f(uow.repos)
	if err != nil {
		return nil, err
	}

	var committedEvents []dorky.Event

	for _, r := range uow.repoList {
		events, err := r.Save()
		if err != nil {
			return nil, err
		}

		committedEvents = append(committedEvents, events...)
	}

	return committedEvents, nil
}
