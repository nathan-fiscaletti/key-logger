package listener

import (
	"context"
	"sync"

	"github.com/nathan-fiscaletti/key-logger/pkg/key"
	"github.com/nathan-fiscaletti/key-logger/pkg/keyboard"
	"github.com/samber/lo"
)

// KeyListener is a listener for key up events, including modifier keys held down at the time of the event.
type KeyListener interface {
	// Events returns a channel of KeyEvents.
	Events(context.Context) (chan KeyEvent, error)
	// Listen listens for keyboard events without blocking and calls the handler for each event.
	Listen(context.Context, KeyEventHandler) error
}

type keyListenerInternal struct {
	keyboard keyboard.Keyboard

	currentModifiers     []key.Key
	currentModifiersLock sync.Mutex

	eventChan chan KeyEvent
}

// New returns a new KeyListener based on the current platform.
func New(kb keyboard.Keyboard) KeyListener {
	return &keyListenerInternal{kb, []key.Key{}, sync.Mutex{}, nil}
}

func (k *keyListenerInternal) Events(ctx context.Context) (chan KeyEvent, error) {
	if k.eventChan == nil {
		k.eventChan = make(chan KeyEvent)

		keyboardEventChan, err := k.keyboard.Events(ctx)
		if err != nil {
			return nil, err
		}

		go func() {
			for {
				select {
				case <-ctx.Done():
					close(k.eventChan)
					k.eventChan = nil
					return
				case event, ok := <-keyboardEventChan:
					if !ok {
						close(k.eventChan)
						k.eventChan = nil
						return
					}

					k.currentModifiersLock.Lock()
					switch event.EventType {
					case keyboard.KeyboardEventTypeDown:
						k.currentModifiers = lo.Uniq(append(k.currentModifiers, event.Key))
					case keyboard.KeyboardEventTypeUp:
						k.currentModifiers = lo.Filter(
							k.currentModifiers,
							func(code key.Key, _ int) bool {
								return code != event.Key
							},
						)
						// Make a copy of the modifiers slice
						k.eventChan <- KeyEvent{
							Key:       event.Key,
							Modifiers: lo.Slice(k.currentModifiers, 0, len(k.currentModifiers)),
						}
					}
					k.currentModifiersLock.Unlock()
				}
			}
		}()
	}

	return k.eventChan, nil
}

func (k *keyListenerInternal) Listen(ctx context.Context, handler KeyEventHandler) error {
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
