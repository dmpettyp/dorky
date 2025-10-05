package dorky

import (
	"time"

	"github.com/dmpettyp/id"
)

type EventID struct{ id.ID }

var NewEventID, MustNewEventID, ParseEventID = id.Intitalizers(
	func(id id.ID) EventID { return EventID{ID: id} },
)

// Event defines the interface required by all dorky domain events. Much of
// this interface is provided when embedding a BaseEvent into the event
// implementation, however EntityID() and EntityType() must be provided by the
// implementation.
type Event interface {
	isEvent()
	isInitialized() bool
	SetEntity(entityType string, entityID id.ID)
	// probalby need Getters eventually
}

// BaseEvent provides an implementation of much of the Event interface which
// can be embedded in specific domain Events defined within client applications
type BaseEvent struct {
	ID          EventID
	Type        string
	Timestamp   time.Time
	EntityType  string
	EntityID    id.ID
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

func (e *BaseEvent) isInitialized() bool {
	return e.initialized
}
