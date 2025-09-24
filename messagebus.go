package dorky

import (
	"context"
)

type MessageBus struct {
}

func NewMessageBus() *MessageBus {
	return &MessageBus{}
}

func (mb *MessageBus) Start(ctx context.Context) {
}

func (mb *MessageBus) HandleCommand(
	ctx context.Context,
	command Command,
) error {
	return nil
}

func (mb *MessageBus) RegisterHandlerMethods(
	ctx context.Context,
	handler any,
) error {
	return nil
}
