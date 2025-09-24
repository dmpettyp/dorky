package dorky

// Entity defines a base dorky Entity type that can be embedded in
// application-defined domain models to get the ability to track domain
// events
type Entity struct {
	Events []Event
}

// AddEvent appends an event to the entity's list of domain events if it has
// been properly initialized
func (entity *Entity) AddEvent(event Event) {
	if event == nil {
		return
	}

	if event.ID().IsNil() {
		return
	}

	entity.Events = append(entity.Events, event)
}

// GetEvents returns the current collection of domain events that have been
// added to the entity
func (entity *Entity) GetEvents() []Event {
	return entity.Events
}

// ResetEvents resets the list of domain events being tracked by the entity
func (entity *Entity) ResetEvents() {
	entity.Events = nil
}
