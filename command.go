package dorky

// Command defines the interface that all command processed by the MessageBus
// must implement
type Command interface {
	isCommand()
}
