package dorky

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sync/atomic"
)

// The MessageBus is a dispatcher for Events and Commands.
//
// It enables clients to register handler methods that are invoked with
// the events or commands that the message bus is processing.
type MessageBus struct {
	started         atomic.Bool
	commandChannel  chan commandWrapper
	eventHandlers   map[reflect.Type][]reflect.Value
	commandHandlers map[reflect.Type]reflect.Value
	eventsToProcess *Queue[Event]
	logger          *slog.Logger
}

func NewMessageBus(logger *slog.Logger) *MessageBus {
	logger.Info("creating MessageBus")

	mb := &MessageBus{
		commandChannel:  make(chan commandWrapper),
		eventHandlers:   make(map[reflect.Type][]reflect.Value),
		commandHandlers: make(map[reflect.Type]reflect.Value),
		eventsToProcess: NewQueue[Event](),
		logger:          logger,
	}

	logger.Info("MessageBus created")

	return mb
}

type commandWrapper struct {
	command Command
	ctx     context.Context
	result  chan error
}

func (mb *MessageBus) Start(ctx context.Context) {
	if !mb.started.CompareAndSwap(false, true) {
		mb.logger.Error("MessageBus already started")
		return
	}

	for {
		select {
		case c := <-mb.commandChannel:
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

func (mb *MessageBus) HandleCommand(
	ctx context.Context,
	command Command,
) error {
	resultChannel := make(chan error)
	defer func() { close(resultChannel) }()

	select {
	case mb.commandChannel <- commandWrapper{command, ctx, resultChannel}:
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

// RegisterHandlerMethods registers each method for the provided handler that
// is a valid handler.
//
// Only one handler for each concrete Command type can be registered, but any
// number of handlers for a concrete Event type can be registered
//
// RegisterHandlerMethods must be called before Start(). It will return an error
// if called after the MessageBus has been started.
func (mb *MessageBus) RegisterHandlerMethods(handler any) error {
	if mb.started.Load() {
		return fmt.Errorf("cannot register handlers after MessageBus has started")
	}

	handlerValue := reflect.ValueOf(handler)

	// Consider registering each handler method
	for i := 0; i < handlerValue.NumMethod(); i++ {
		handlerName := fmt.Sprintf(
			"%T::%s", handler, handlerValue.Type().Method(i).Name,
		)

		mb.logger.Info("registering handler method", "name", handlerName)

		handlerFunc := handlerValue.Method(i)

		handlerArgType, handlerArgKind, err := validateHandler(handlerFunc)

		if err != nil {
			// Skip methods that don't match handler signature - they're not handlers
			mb.logger.Debug("skipping non-handler method", "name", handlerName, "reason", err)
			continue
		}

		if handlerArgKind == handlerArgKindCommand {
			if _, ok := mb.commandHandlers[handlerArgType]; ok {
				return fmt.Errorf(
					"duplicate command handler for %v: %s (handler already registered)",
					handlerArgType, handlerName,
				)
			}

			mb.commandHandlers[handlerArgType] = handlerFunc

			mb.logger.Info("registered command handler", "name", handlerName)
		} else if handlerArgKind == handlerArgKindEvent {
			mb.eventHandlers[handlerArgType] = append(
				mb.eventHandlers[handlerArgType],
				handlerFunc,
			)

			mb.logger.Info("registered event handler", "name", handlerName)
		}
	}

	return nil
}

func (mb *MessageBus) dispatchCommand(ctx context.Context, command Command) error {
	mb.logger.Info(
		"messagebus dispatching command",
		// "messageEntity", comand.GetEntity(),
		// "messageType", message.GetType(),
		// "message", messageToString(message),
	)

	commandType := reflect.TypeOf(command)

	handler, ok := mb.commandHandlers[commandType]

	if !ok {
		mb.logger.Info("no command handler found")
		return fmt.Errorf("no handler for command")
	}

	mb.logger.Info("invoking command handler")

	events, err := invokeHandler(ctx, handler, command)

	if err != nil {
		mb.logger.Info(
			"invoking command handler failed",
			"error", err.Error(),
		)
		return err
	}

	mb.eventsToProcess.enqueueMultiple(events)

	return nil
}

func (mb *MessageBus) dispatchEvents(ctx context.Context) {
	for {
		event, ok := mb.eventsToProcess.dequeue()

		if !ok {
			return
		}

		mb.logger.Info(
			"messagebus dispatching event",
			// "messageEntity", message.GetEntity(),
			// "messageType", message.GetType(),
			// "message", messageToString(message),
		)

		messageType := reflect.TypeOf(event)

		if eventHandlers, ok := mb.eventHandlers[messageType]; ok {
			for _, handler := range eventHandlers {
				mb.logger.Info("invoking event handler")

				events, err := invokeHandler(ctx, handler, event)

				if err != nil {
					mb.logger.Info(
						"invoking event handler failed",
						"error", err.Error(),
					)
				}

				mb.eventsToProcess.enqueueMultiple(events)
			}
		}
	}
}
