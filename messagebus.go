package dorky

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sync"
	"sync/atomic"
)

// The MessageBus is a dispatcher for Events and Commands.
//
// It enables clients to register handler methods that are invoked with
// the events or commands that the message bus is processing.
type MessageBus struct {
	started         atomic.Bool
	commands        chan messageBusCommand
	eventHandlers   map[reflect.Type][]func(context.Context, Event) ([]Event, error)
	commandHandlers map[reflect.Type]func(context.Context, Command) ([]Event, error)
	eventsToProcess *Queue[Event]
	wg              sync.WaitGroup
	logger          *slog.Logger
}

func NewMessageBus(logger *slog.Logger) *MessageBus {
	logger.Info("creating MessageBus")

	mb := &MessageBus{
		commands:        make(chan messageBusCommand),
		eventHandlers:   make(map[reflect.Type][]func(context.Context, Event) ([]Event, error)),
		commandHandlers: make(map[reflect.Type]func(context.Context, Command) ([]Event, error)),
		eventsToProcess: NewQueue[Event](),
		logger:          logger,
	}

	logger.Info("MessageBus created")

	return mb
}

// RegisterCommandHandler registers a type-safe command handler with the
// MessageBus. This is implemented as a function that calls
// registerCommandHandler on the MessageBus because generic methods are not
// allowed.
func RegisterCommandHandler[C Command](
	mb *MessageBus,
	handler func(context.Context, C) ([]Event, error),
) error {
	var zero C

	return mb.registerCommandHandler(
		reflect.TypeOf(zero),
		func(ctx context.Context, cmd Command) ([]Event, error) {
			return handler(ctx, cmd.(C))
		},
	)
}

// RegisterEvent registers a type-safe event handler with the MessageBus. This
// is implemented as a function that calls registerEventHandler on the
// MessageBus because generic methods are not allowed.
func RegisterEventHandler[E Event](
	mb *MessageBus,
	handler func(context.Context, E) ([]Event, error),
) error {
	var zero E

	return mb.registerEventHandler(
		reflect.TypeOf(zero),
		func(ctx context.Context, evt Event) ([]Event, error) {
			return handler(ctx, evt.(E))
		},
	)
}

type messageBusCommand struct {
	command Command
	ctx     context.Context
	result  chan error
}

func (mb *MessageBus) Start(ctx context.Context) {
	if !mb.started.CompareAndSwap(false, true) {
		mb.logger.Error("MessageBus already started")
		return
	}

	mb.wg.Add(1)

	defer mb.wg.Done()

	mb.logger.Info("starting MessageBus")

	for {
		select {
		case c, ok := <-mb.commands:
			if !ok {
				// Channel closed, shutdown
				return
			}

			err := mb.dispatchCommand(c.ctx, c.command)

			select {
			case c.result <- err:
			case <-c.ctx.Done():
			}

			mb.dispatchEvents(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (mb *MessageBus) Stop() {
	mb.logger.Info("stopping MessageBus")
	close(mb.commands)
	mb.wg.Wait()
}

func (mb *MessageBus) HandleCommand(
	ctx context.Context,
	command Command,
) error {
	resultChannel := make(chan error)
	defer func() { close(resultChannel) }()

	select {
	case mb.commands <- messageBusCommand{command, ctx, resultChannel}:
	case <-ctx.Done():
		return fmt.Errorf(
			"cannot send a command to the messagebus to handle: %w", ctx.Err(),
		)
	}

	select {
	case result := <-resultChannel:
		return result
	case <-ctx.Done():
		return fmt.Errorf(
			"cannot receive messagebus handle response: %w", ctx.Err(),
		)
	}
}

// registerCommandHandler registers a type safe handler for the commandType
// provided. Only one handler may be registered for each commandType
func (mb *MessageBus) registerCommandHandler(
	commandType reflect.Type,
	handler func(context.Context, Command) ([]Event, error),
) error {
	if mb.started.Load() {
		return fmt.Errorf("cannot register handlers after MessageBus has started")
	}

	if _, exists := mb.commandHandlers[commandType]; exists {
		return fmt.Errorf("handler already registered for command type %v", commandType)
	}

	mb.commandHandlers[commandType] = handler

	mb.logger.Info("registered command handler", "type", commandType)

	return nil
}

// registerEventHandler registers a type safe handler for the Event type
// provided. Many handler may be registered for each Event type
func (mb *MessageBus) registerEventHandler(
	eventType reflect.Type,
	handler func(context.Context, Event) ([]Event, error),
) error {
	if mb.started.Load() {
		return fmt.Errorf("cannot register event handler after MessageBus has started")
	}

	mb.eventHandlers[eventType] = append(
		mb.eventHandlers[eventType],
		handler,
	)

	mb.logger.Info("registered event handler", "type", eventType)

	return nil
}

// dispatchCommand invokes the command handler for the type of Command
// passed in. Events generated from invoking the handler are queued and
// dispatched to event handlers after the command handler returns.
func (mb *MessageBus) dispatchCommand(ctx context.Context, command Command) error {
	mb.logger.Info("messagebus dispatching command")

	commandType := reflect.TypeOf(command)

	handler, ok := mb.commandHandlers[commandType]

	if !ok {
		mb.logger.Info("no command handler found")
		return fmt.Errorf("no handler for command type %v", commandType)
	}

	mb.logger.Info("invoking command handler")

	events, err := handler(ctx, command)

	if err != nil {
		mb.logger.Info("invoking command handler failed", "error", err.Error())
		return err
	}

	mb.eventsToProcess.enqueueMultiple(events)

	return nil
}

// dispatchEvents dispatches all events queued up by the MessageBus to any
// handlers that are registered for them. Events returned by the event
// handlers are queued up and processed before returning.
func (mb *MessageBus) dispatchEvents(ctx context.Context) {
	for {
		event, ok := mb.eventsToProcess.dequeue()

		if !ok {
			return
		}

		mb.logger.Info("messagebus dispatching event")

		eventType := reflect.TypeOf(event)

		if handlers, ok := mb.eventHandlers[eventType]; ok {
			for _, handler := range handlers {
				mb.logger.Info("invoking event handler")

				events, err := handler(ctx, event)

				if err != nil {
					mb.logger.Info("invoking event handler failed", "error", err.Error())
				}

				mb.eventsToProcess.enqueueMultiple(events)
			}
		}
	}
}
