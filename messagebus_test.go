package dorky_test

import (
	"context"
	//  "io"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dmpettyp/dorky"
)

// var logger = slog.New(slog.NewTextHandler(io.Discard, nil))
var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

// Create a new Logger using the handler

// func TestInvalidHandlers(t *testing.T) {
// 	for _, tc := range []struct {
// 		name    string
// 		handler any
// 	}{
// 		{"struct as handler", struct{}{}},
// 		{
// 			name:    "string as handler",
// 			handler: "i'm a handler",
// 		},
// 		{
// 			name:    "handler doesn't accept any arguments",
// 			handler: func() ([]ddd.Message, error) { return nil, nil },
// 		},
// 		{
// 			name:    "handler doesn't return anything",
// 			handler: func(int) {},
// 		},
// 		{
// 			name:    "handler only returns error",
// 			handler: func(int) error { return nil },
// 		},
// 		{
// 			name: "handler invalid events returned",
// 			handler: func(int) (int, error) {
// 				return 0, nil
// 			},
// 		},
// 		{
// 			name: "handler doesn't return Message slice",
// 			handler: func(int) ([]int, error) {
// 				return nil, nil
// 			},
// 		},
// 		{
// 			name: "handler doesn't return error",
// 			handler: func(int) ([]ddd.Message, int) {
// 				return nil, 0
// 			},
// 		},
// 		{
// 			name: "handle doesn't have ctx arg",
// 			handler: func(int) ([]ddd.Message, error) {
// 				return nil, nil
// 			},
// 		},
// 	} {
// 		t.Run(tc.name, func(t *testing.T) {
// 			mb := ddd.CreateMessageBus()
// 			err := mb.RegisterEventHandler(tc.handler)
// 			assert.Error(t, err)
// 		})
// 	}
// }

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

// func TestValidHandlers(t *testing.T) {
// 	mb := ddd.CreateMessageBus()
//
// 	eventHandler := &eventHandler1{}
// 	commandHandler := &commandHandler1{}
//
// 	err := mb.RegisterEventHandler(eventHandler.Handle)
// 	assert.NoError(t, err)
//
// 	err = mb.RegisterCommandHandler(commandHandler.Handle)
// 	assert.NoError(t, err)
// }

// func TestMultipleEventHandlersValid(t *testing.T) {
// 	mb := ddd.CreateMessageBus()
//
// 	handler1 := &eventHandler1{}
// 	handler2 := &eventHandler1{}
//
// 	err := mb.RegisterEventHandler(handler1.Handle)
// 	assert.NoError(t, err)
//
// 	err = mb.RegisterEventHandler(handler2.Handle)
// 	assert.NoError(t, err)
// }

// func TestCantRegisterForEventAndCommandHandler(t *testing.T) {
// 	mb := ddd.CreateMessageBus()
//
// 	handler1 := &eventHandler1{}
// 	handler2 := &eventHandler1{}
//
// 	err := mb.RegisterEventHandler(handler1.Handle)
// 	assert.NoError(t, err)
//
// 	err = mb.RegisterCommandHandler(handler2.Handle)
// 	assert.Error(t, err)
// }
//
// func TestCantRegisterForCommandAndEventHandler(t *testing.T) {
// 	mb := ddd.CreateMessageBus()
//
// 	handler1 := &eventHandler1{}
// 	handler2 := &eventHandler1{}
//
// 	err := mb.RegisterCommandHandler(handler2.Handle)
// 	assert.NoError(t, err)
//
// 	err = mb.RegisterEventHandler(handler1.Handle)
// 	assert.Error(t, err)
// }
//
// func TestMultipleCommandHandlersInvalid(t *testing.T) {
// 	mb := ddd.CreateMessageBus()
//
// 	handler1 := &commandHandler1{}
// 	handler2 := &commandHandler1{}
//
// 	err := mb.RegisterCommandHandler(handler1.Handle)
// 	assert.NoError(t, err)
//
// 	err = mb.RegisterCommandHandler(handler2.Handle)
// 	assert.Error(t, err)
// }
//
// func TestDifferentCommandHandlersValid(t *testing.T) {
// 	mb := ddd.CreateMessageBus()
//
// 	handler1 := &commandHandler1{}
// 	handler2 := &commandHandler2{}
//
// 	err := mb.RegisterCommandHandler(handler1.Handle)
// 	assert.NoError(t, err)
//
// 	err = mb.RegisterCommandHandler(handler2.Handle)
// 	assert.NoError(t, err)
// }
//
// func TestCommandHandlersInvoked(t *testing.T) {
// 	mb := ddd.CreateMessageBus()
//
// 	handler1 := &commandHandler1{}
// 	handler2 := &commandHandler2{}
//
// 	err := mb.RegisterCommandHandler(handler1.Handle)
// 	assert.NoError(t, err)
//
// 	err = mb.RegisterCommandHandler(handler2.Handle)
// 	assert.NoError(t, err)
//
// 	go mb.Start(context.Background())
//
// 	ctx := context.Background()
//
// 	m1 := commandArg1{value: "string message"}
// 	_ = mb.Handle(ctx, &m1)
//
// 	assert.Equal(t, 1, handler1.CallCount)
// 	assert.Equal(t, 0, handler2.CallCount)
//
// 	m2 := commandArg2{value: "string message"}
// 	_ = mb.Handle(ctx, &m2)
//
// 	assert.Equal(t, 1, handler1.CallCount)
// 	assert.Equal(t, 1, handler2.CallCount)
// }
//
// func TestEventHandlersInvoked(t *testing.T) {
// 	mb := ddd.CreateMessageBus()
//
// 	handler1 := &eventHandler1{}
// 	handler2 := &eventHandler1{}
// 	handler3 := &eventHandler2{}
//
// 	err := mb.RegisterEventHandler(handler1.Handle)
// 	assert.NoError(t, err)
//
// 	err = mb.RegisterEventHandler(handler2.Handle)
// 	assert.NoError(t, err)
//
// 	err = mb.RegisterEventHandler(handler3.Handle)
// 	assert.NoError(t, err)
//
// 	go mb.Start(context.Background())
//
// 	m1 := eventArg1{value: "event details"}
//
// 	ctx := context.Background()
//
// 	_ = mb.Handle(ctx, &m1)
//
// 	assert.Equal(t, 1, handler1.CallCount)
// 	assert.Equal(t, 1, handler2.CallCount)
// 	assert.Equal(t, 0, handler3.CallCount)
//
// 	_ = mb.Handle(ctx, &m1)
//
// 	assert.Equal(t, 2, handler1.CallCount)
// 	assert.Equal(t, 2, handler2.CallCount)
// 	assert.Equal(t, 0, handler3.CallCount)
//
// 	m2 := eventArg2{value: "event details"}
//
// 	_ = mb.Handle(ctx, &m2)
//
// 	assert.Equal(t, 2, handler1.CallCount)
// 	assert.Equal(t, 2, handler2.CallCount)
// 	assert.Equal(t, 1, handler3.CallCount)
// }
//
// func TestHandleErrorsIfHandlerContextExpires(t *testing.T) {
// 	mb := ddd.CreateMessageBus()
//
// 	handler := &eventHandler1{}
// 	err := mb.RegisterEventHandler(handler.Handle)
// 	assert.NoError(t, err)
//
// 	// Don't start the messagebus so messages aren't handled and the
// 	// context deadline occurs
// 	// go mb.Start(context.Background())
//
// 	m1 := eventArg1{value: "event details"}
//
// 	d := time.Now().Add(1 * time.Millisecond)
// 	ctx, cancel := context.WithDeadline(context.Background(), d)
//
// 	err = mb.Handle(ctx, &m1)
// 	assert.Error(t, err)
//
// 	assert.Equal(t, 0, handler.CallCount)
// 	cancel()
// }
//
// type orderEvent struct {
// 	ddd.BaseMessage
// 	id int
// }
//
// type orderEventHandler struct {
// 	expected int
// 	t        *testing.T
// }
//
// func (h *orderEventHandler) Handle(
// 	_ context.Context,
// 	event *orderEvent,
// ) (
// 	[]ddd.Message,
// 	error,
// ) {
// 	assert.Equal(h.t, h.expected, event.id)
//
// 	h.expected++
//
// 	if event.id == 1 {
// 		return []ddd.Message{
// 			&orderEvent{id: 2},
// 			&orderEvent{id: 3},
// 			&orderEvent{id: 4},
// 		}, nil
// 	}
//
// 	if event.id == 2 {
// 		return []ddd.Message{
// 			&orderEvent{id: 5},
// 			&orderEvent{id: 6},
// 			&orderEvent{id: 7},
// 		}, nil
// 	}
//
// 	if event.id == 6 {
// 		return []ddd.Message{
// 			&orderEvent{id: 8},
// 		}, nil
// 	}
//
// 	if event.id == 8 {
// 		return []ddd.Message{
// 			&orderEvent{id: 9},
// 		}, nil
// 	}
//
// 	return nil, nil
// }
//
// func TestDispatchOrder(t *testing.T) {
// 	mb := ddd.CreateMessageBus()
//
// 	handler := &orderEventHandler{expected: 1, t: t}
// 	err := mb.RegisterEventHandler(handler.Handle)
// 	assert.NoError(t, err)
//
// 	go mb.Start(context.Background())
//
// 	_ = mb.Handle(context.Background(), &orderEvent{id: 1})
//
// 	assert.Equal(t, 10, handler.expected)
// }

// type event1 struct {
// 	ddd.BaseEvent
// }
//
// type command1 struct {
// 	ddd.BaseCommand
// }
//
// type foo struct {
// 	eventHandlerInvocations   int
// 	commandHandlerInvocations int
// }
//
// func (f *foo) ThisIsAnEventHandler(
// 	_ context.Context,
// 	_ *event1,
// ) ([]ddd.Message, error) {
// 	f.eventHandlerInvocations += 1
// 	return nil, nil
// }
//
// func (f *foo) ThisIsACommandHandler(
// 	_ context.Context,
// 	_ *command1,
// ) ([]ddd.Message, error) {
// 	f.commandHandlerInvocations += 1
// 	return nil, nil
// }

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

	require.Equal(t, 1, cmd1Count)
	require.Equal(t, 1, cmd2Count)
	require.Equal(t, 1, evt1Count)
	require.Equal(t, 1, evt2Count)
}
