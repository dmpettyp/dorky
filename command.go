package dorky

// Command defines the interface that all command processed by the MessageBus
// must implement
type Command interface {
	isCommand()
}

type BaseCommand struct {
}

// BaseCommand must implement isCommand to be recognized as a dorky Command
func (*BaseCommand) isCommand() {}
