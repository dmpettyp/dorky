package messages

import (
	"time"

	"github.com/dmpettyp/dorky/id"
)

type EventID struct{ id.ID }

var NewEventID, MustNewEventID, ParseEventID = id.Create(
	func(id id.ID) EventID { return EventID{ID: id} },
)

// Event defines the interface required by all dorky domain events. Much of
// this interface is provided when embedding a BaseEvent into the event
// implementation, however EntityID() and EntityType() must be provided by the
// implementation.
type Event interface {
	isEvent()
	IsInitialized() bool
	SetEntity(entityType string, entityID id.ID)
	GetType() string
	GetTimestamp() time.Time
	GetEntityID() id.ID
	GetEntityType() string
}

// BaseEvent provides an implementation of much of the Event interface which
// can be embedded in specific domain Events defined within client applications
type BaseEvent struct {
	ID          EventID   `json:"id"`
	Type        string    `json:"type"`
	Timestamp   time.Time `json:"timestamp"`
	EntityType  string    `json:"entity_type"`
	EntityID    id.ID     `json:"entity_id"`
	initialized bool
}

// BaseEvent must implement isEvent to be recognized as a dorky Event
func (*BaseEvent) isEvent() {}

// Init sets the BaseEvent eventType, entityType and entityID, and initializes
// its ID and timestamp
func (e *BaseEvent) Init(eventType string) {
	e.ID, _ = NewEventID()
	e.Timestamp = time.Now().UTC()
	e.Type = eventType
	e.initialized = true
}

func (e *BaseEvent) SetEntity(entityType string, entityID id.ID) {
	e.EntityType = entityType
	e.EntityID = entityID
}

func (e *BaseEvent) IsInitialized() bool {
	return e.initialized
}

func (e *BaseEvent) GetType() string {
	return e.Type
}

func (e *BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e *BaseEvent) GetEntityID() id.ID {
	return e.EntityID
}

func (e *BaseEvent) GetEntityType() string {
	return e.EntityType
}
