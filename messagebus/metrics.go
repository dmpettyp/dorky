package messagebus

import (
	"log/slog"
	"time"
)

type MetricsHook interface {
	ObserveCommand(commandType string, status string, duration time.Duration)
	ObserveEvent(eventType string, status string, duration time.Duration)
}

type Option func(*MessageBus)

func WithMetricsHook(hook MetricsHook) Option {
	return func(mb *MessageBus) {
		mb.metrics = hook
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(mb *MessageBus) {
		mb.logger = logger
	}
}
