package dorky

import (
// "encoding/json"
// "fmt"
)

// Entity defines a base dorky Entity type that can be embedded in
// application-defined domain models to get the ability to track domain
// events
type Entity struct {
	Events []Event
}

// AddEvent appends an event to the entity's list of domain events if it has
// been properly initialized
func (entity *Entity) AddEvent(e Event) {
	if e == nil {
		return
	}

	if !e.isInitialized() {
		return
	}

	entity.Events = append(entity.Events, e)

	// jsonData, err := json.Marshal(e)
	//
	// if err != nil {
	// 	return
	// }
	//
	// fmt.Println(string(jsonData))
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
