package dorky_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	// "os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dmpettyp/dorky"
)

var logger = slog.New(slog.NewTextHandler(io.Discard, nil))

// var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

// Define valid arguments for event and command handlers
type commandArg struct {
	dorky.BaseCommand
	value string
}

type commandArg1 commandArg
type commandArg2 commandArg

type eventArg struct {
	dorky.BaseEvent
	value string
}

type eventArg1 eventArg
type eventArg2 eventArg

func TestTypeSafeRegistration(t *testing.T) {
	mb := dorky.NewMessageBus(logger)

	cmd1Count := 0
	cmd2Count := 0
	evt1Count := 0
	evt2Count := 0

	// Type-safe command handlers
	err := dorky.RegisterCommandHandler(mb, func(ctx context.Context, cmd *commandArg1) ([]dorky.Event, error) {
		cmd1Count++
		return []dorky.Event{&eventArg1{}, &eventArg2{}}, nil
	})
	require.NoError(t, err)

	err = dorky.RegisterCommandHandler(mb, func(ctx context.Context, cmd *commandArg2) ([]dorky.Event, error) {
		cmd2Count++
		return nil, nil
	})
	require.NoError(t, err)

	// Type-safe event handlers
	err = dorky.RegisterEventHandler(mb, func(ctx context.Context, evt *eventArg1) ([]dorky.Event, error) {
		evt1Count++
		return nil, nil
	})
	require.NoError(t, err)

	err = dorky.RegisterEventHandler(mb, func(ctx context.Context, evt *eventArg2) ([]dorky.Event, error) {
		evt2Count++
		return nil, nil
	})
	require.NoError(t, err)

	go mb.Start(context.Background())

	_ = mb.HandleCommand(context.Background(), &commandArg1{})
	_ = mb.HandleCommand(context.Background(), &commandArg2{})

	mb.Stop()

	require.Equal(t, 1, cmd1Count)
	require.Equal(t, 1, cmd2Count)
	require.Equal(t, 1, evt1Count)
	require.Equal(t, 1, evt2Count)
}

// Types for method-based handler test
type CreateOrderCommand struct {
	dorky.BaseCommand
	OrderID string
}

type OrderCreatedEvent struct {
	dorky.BaseEvent
	OrderID string
}

type ShipOrderCommand struct {
	dorky.BaseCommand
	OrderID string
}

type OrderShippedEvent struct {
	dorky.BaseEvent
	OrderID string
}

type OrderService struct {
	ordersCreated int
	ordersShipped int
}

func (s *OrderService) HandleCreateOrder(ctx context.Context, cmd *CreateOrderCommand) ([]dorky.Event, error) {
	return []dorky.Event{&OrderCreatedEvent{OrderID: cmd.OrderID}}, nil
}

func (s *OrderService) HandleShipOrder(ctx context.Context, cmd *ShipOrderCommand) ([]dorky.Event, error) {
	return []dorky.Event{&OrderShippedEvent{OrderID: cmd.OrderID}}, nil
}

func (s *OrderService) OnOrderCreated(ctx context.Context, evt *OrderCreatedEvent) ([]dorky.Event, error) {
	s.ordersCreated++
	return nil, nil
}

func (s *OrderService) OnOrderShipped(ctx context.Context, evt *OrderShippedEvent) ([]dorky.Event, error) {
	s.ordersShipped++
	return nil, nil
}

// Test that method-based handlers work (handlers as methods on structs)
func TestMethodBasedHandlers(t *testing.T) {
	mb := dorky.NewMessageBus(logger)

	svc := &OrderService{}

	// Register methods as handlers
	err := dorky.RegisterCommandHandler(mb, svc.HandleCreateOrder)
	require.NoError(t, err)

	err = dorky.RegisterCommandHandler(mb, svc.HandleShipOrder)
	require.NoError(t, err)

	err = dorky.RegisterEventHandler(mb, svc.OnOrderCreated)
	require.NoError(t, err)

	err = dorky.RegisterEventHandler(mb, svc.OnOrderShipped)
	require.NoError(t, err)

	go mb.Start(context.Background())

	// Execute commands
	err = mb.HandleCommand(context.Background(), &CreateOrderCommand{OrderID: "123"})
	require.NoError(t, err)

	err = mb.HandleCommand(context.Background(), &ShipOrderCommand{OrderID: "123"})
	require.NoError(t, err)

	mb.Stop()

	// Verify state was updated through methods
	require.Equal(t, 1, svc.ordersCreated)
	require.Equal(t, 1, svc.ordersShipped)
}

// Test event cascade - events generating more events (breadth-first)
func TestEventCascade(t *testing.T) {
	mb := dorky.NewMessageBus(logger)

	type TriggerCommand struct {
		dorky.BaseCommand
	}

	type Level1Event struct {
		dorky.BaseEvent
		ID int
	}

	type Level2Event struct {
		dorky.BaseEvent
		ID int
	}

	type Level3Event struct {
		dorky.BaseEvent
		ID int
	}

	var executionOrder []string

	// Command generates 2 level-1 events
	err := dorky.RegisterCommandHandler(mb, func(ctx context.Context, cmd *TriggerCommand) ([]dorky.Event, error) {
		executionOrder = append(executionOrder, "command")
		return []dorky.Event{&Level1Event{ID: 1}, &Level1Event{ID: 2}}, nil
	})
	require.NoError(t, err)

	// Each level-1 event generates a level-2 event
	err = dorky.RegisterEventHandler(mb, func(ctx context.Context, evt *Level1Event) ([]dorky.Event, error) {
		executionOrder = append(executionOrder, fmt.Sprintf("level1-%d", evt.ID))
		return []dorky.Event{&Level2Event{ID: evt.ID}}, nil
	})
	require.NoError(t, err)

	// Each level-2 event generates a level-3 event
	err = dorky.RegisterEventHandler(mb, func(ctx context.Context, evt *Level2Event) ([]dorky.Event, error) {
		executionOrder = append(executionOrder, fmt.Sprintf("level2-%d", evt.ID))
		return []dorky.Event{&Level3Event{ID: evt.ID}}, nil
	})
	require.NoError(t, err)

	// Level-3 events don't generate more events
	err = dorky.RegisterEventHandler(mb, func(ctx context.Context, evt *Level3Event) ([]dorky.Event, error) {
		executionOrder = append(executionOrder, fmt.Sprintf("level3-%d", evt.ID))
		return nil, nil
	})
	require.NoError(t, err)

	go mb.Start(context.Background())

	err = mb.HandleCommand(context.Background(), &TriggerCommand{})
	require.NoError(t, err)

	mb.Stop()

	// Verify breadth-first order: command, then all level-1, then all level-2, then all level-3
	expected := []string{
		"command",
		"level1-1", "level1-2",
		"level2-1", "level2-2",
		"level3-1", "level3-2",
	}
	require.Equal(t, expected, executionOrder)
}

// Test that command handler errors propagate correctly and events aren't dispatched
func TestCommandHandlerErrors(t *testing.T) {
	mb := dorky.NewMessageBus(logger)

	type FailingCommand struct {
		dorky.BaseCommand
	}

	type SuccessEvent struct {
		dorky.BaseEvent
	}

	eventHandlerCalled := false

	// Command handler that returns an error
	err := dorky.RegisterCommandHandler(mb, func(ctx context.Context, cmd *FailingCommand) ([]dorky.Event, error) {
		// Return events AND an error - events should NOT be dispatched
		return []dorky.Event{&SuccessEvent{}}, fmt.Errorf("command failed")
	})
	require.NoError(t, err)

	// Event handler should not be called when command fails
	err = dorky.RegisterEventHandler(mb, func(ctx context.Context, evt *SuccessEvent) ([]dorky.Event, error) {
		eventHandlerCalled = true
		return nil, nil
	})
	require.NoError(t, err)

	go mb.Start(context.Background())

	err = mb.HandleCommand(context.Background(), &FailingCommand{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "command failed")

	mb.Stop()

	// Event handler should NOT have been called
	require.False(t, eventHandlerCalled)
}

// Test that event handler errors are logged but don't stop other handlers
func TestEventHandlerErrors(t *testing.T) {
	mb := dorky.NewMessageBus(logger)

	type TriggerCommand struct {
		dorky.BaseCommand
	}

	type ProcessEvent struct {
		dorky.BaseEvent
	}

	handler1Called := false
	handler2Called := false
	handler3Called := false

	// Command generates an event
	err := dorky.RegisterCommandHandler(mb, func(ctx context.Context, cmd *TriggerCommand) ([]dorky.Event, error) {
		return []dorky.Event{&ProcessEvent{}}, nil
	})
	require.NoError(t, err)

	// First handler succeeds
	err = dorky.RegisterEventHandler(mb, func(ctx context.Context, evt *ProcessEvent) ([]dorky.Event, error) {
		handler1Called = true
		return nil, nil
	})
	require.NoError(t, err)

	// Second handler fails
	err = dorky.RegisterEventHandler(mb, func(ctx context.Context, evt *ProcessEvent) ([]dorky.Event, error) {
		handler2Called = true
		return nil, fmt.Errorf("handler 2 failed")
	})
	require.NoError(t, err)

	// Third handler succeeds
	err = dorky.RegisterEventHandler(mb, func(ctx context.Context, evt *ProcessEvent) ([]dorky.Event, error) {
		handler3Called = true
		return nil, nil
	})
	require.NoError(t, err)

	go mb.Start(context.Background())

	err = mb.HandleCommand(context.Background(), &TriggerCommand{})
	require.NoError(t, err) // Command should still succeed

	mb.Stop()

	// All handlers should have been called despite handler 2 failing
	require.True(t, handler1Called)
	require.True(t, handler2Called)
	require.True(t, handler3Called)
}

// Test context cancellation behavior
func TestContextCancellation(t *testing.T) {
	mb := dorky.NewMessageBus(logger)

	type SlowCommand struct {
		dorky.BaseCommand
	}

	err := dorky.RegisterCommandHandler(mb, func(ctx context.Context, cmd *SlowCommand) ([]dorky.Event, error) {
		return nil, nil
	})
	require.NoError(t, err)

	go mb.Start(context.Background())

	// Create an already-cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = mb.HandleCommand(ctx, &SlowCommand{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot send a command to the messagebus to handle")

	mb.Stop()
}

// Test duplicate handler registration behavior
func TestDuplicateHandlerRegistration(t *testing.T) {
	mb := dorky.NewMessageBus(logger)

	type MyCommand struct {
		dorky.BaseCommand
	}

	type MyEvent struct {
		dorky.BaseEvent
	}

	// Register first command handler - should succeed
	err := dorky.RegisterCommandHandler(mb, func(ctx context.Context, cmd *MyCommand) ([]dorky.Event, error) {
		return nil, nil
	})
	require.NoError(t, err)

	// Register duplicate command handler - should fail
	err = dorky.RegisterCommandHandler(mb, func(ctx context.Context, cmd *MyCommand) ([]dorky.Event, error) {
		return nil, nil
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "handler already registered")

	// Register first event handler - should succeed
	err = dorky.RegisterEventHandler(mb, func(ctx context.Context, evt *MyEvent) ([]dorky.Event, error) {
		return nil, nil
	})
	require.NoError(t, err)

	// Register second event handler for same type - should succeed (multiple allowed)
	err = dorky.RegisterEventHandler(mb, func(ctx context.Context, evt *MyEvent) ([]dorky.Event, error) {
		return nil, nil
	})
	require.NoError(t, err)
}
