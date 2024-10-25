package keyboard

import (
	"context"
)

var keyboard EventChannel = nil

type Keyboard interface {
	EventChannel
	// Listen listens for keyboard events without blocking and calls the handler for each event.
	Listen(context.Context, EventHandler) error
}

type EventChannel interface {
	// Events returns a channel of KeyboardEvents.
	Events(context.Context) (chan Event, error)
}

type keyboardInternal struct {
	EventChannel
}

// New returns a new Keyboard based on the current platform.
func New() Keyboard {
	return keyboardInternal{keyboard}
}

func (k keyboardInternal) Listen(ctx context.Context, handler EventHandler) error {
	events, err := k.Events(ctx)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-events:
				handler(event)
			}
		}
	}()

	return nil
}
