package dorky

import (
	"time"

	"github.com/dmpettyp/id"
)

// Event defines the interface required by all dorky domain events. Much of
// this interface is provided when embedding a BaseEvent into the event
// implementation, however EntityID() and EntityType() must be provided by the
// implementation.
type Event interface {
	isEvent()
	ID() EventID
	Type() string
	Timestamp() time.Time
	EntityID() id.ID
	EntityType() string
}

// BaseEvent provides an implementation of much of the Event interface which
// can be embedded in specific domain Events defined within client applications
type BaseEvent struct {
	eventID   EventID
	eventType string
	timestamp time.Time
}

// BaseEvent must implement isEvent to be recognized as a dorky Event
func (*BaseEvent) isEvent() {}

// Init sets the BaseEvent eventType and initializes its ID and timestamp
func (e *BaseEvent) Init(eventType string) {
	e.eventID, _ = NewEventID()
	e.timestamp = time.Now().UTC()
	e.eventType = eventType
}

// ID returns the EventId of the BaseEvent
func (e *BaseEvent) ID() EventID {
	return e.eventID
}

// ID returns the Type of the BaseEvent
func (e *BaseEvent) Type() string {
	return e.eventType
}

// ID returns the Timestamp of the BaseEvent
func (e *BaseEvent) Timestamp() time.Time {
	return e.timestamp
}
