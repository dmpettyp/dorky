package aggregate

import (
	"github.com/dmpettyp/dorky/messages"
)

// Aggregate defines a base dorky Aggregate type that can be embedded in
// application-defined domain models to get the ability to track domain
// events
type Aggregate struct {
	Events []messages.Event
}

// AddEvent appends an event to the entity's list of domain events if it has
// been properly initialized
func (aggregate *Aggregate) AddEvent(event messages.Event) {
	if event == nil {
		return
	}

	if !event.IsInitialized() {
		return
	}

	aggregate.Events = append(aggregate.Events, event)
}

// GetEvents returns the current collection of domain events that have been
// added to the entity
func (aggregate *Aggregate) GetEvents() []messages.Event {
	return aggregate.Events
}

// ResetEvents resets the list of domain events being tracked by the entity
func (aggregate *Aggregate) ResetEvents() {
	aggregate.Events = nil
}
