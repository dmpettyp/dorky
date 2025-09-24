package dorky

import "github.com/dmpettyp/id"

type EventID struct{ id.ID }

var NewEventID, MustNewEventID, ParseEventID = id.Intitalizers(
	func(id id.ID) EventID { return EventID{ID: id} },
)
