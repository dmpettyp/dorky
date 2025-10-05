package dorky

import (
	"time"

	"github.com/dmpettyp/id"
)

type CommandID struct{ id.ID }

var NewCommandID, MustNewCommandID, ParseCommandID = id.Inititalizers(
	func(id id.ID) CommandID { return CommandID{ID: id} },
)

// Command defines the interface that all command processed by the MessageBus
// must implement
type Command interface {
	isCommand()
}

type BaseCommand struct {
	ID        CommandID
	Type      string
	Timestamp time.Time
}

// BaseCommand must implement isCommand to be recognized as a dorky Command
func (*BaseCommand) isCommand() {}

func (c *BaseCommand) Init(commandType string) {
	c.ID, _ = NewCommandID()
	c.Timestamp = time.Now().UTC()
	c.Type = commandType
}
