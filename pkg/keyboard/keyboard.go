package keyboard

import (
	"context"
	"sync"
	"time"

	"github.com/nathan-fiscaletti/key-logger/pkg/keyboard/key"
	"github.com/samber/lo"
)

type KeyboardEventType string

const (
	KeyboardEventTypeUp   KeyboardEventType = "up"
	KeyboardEventTypeDown KeyboardEventType = "down"
)

type KeyboardEvent struct {
	Key       key.Key
	EventType KeyboardEventType
	Timestamp time.Time
}

var keyboard Keyboard = nil

type Keyboard interface {
	Events(context.Context) (chan KeyboardEvent, error)
}

func NewKeyboard() Keyboard {
	return keyboard
}

type KeyEvent struct {
	Modifiers []key.Key
	Key       key.Key
}

type KeyListener struct {
	keyboard Keyboard

	currentModifiers     []key.Key
	currentModifiersLock sync.Mutex

	eventChan chan KeyEvent
}

func NewKeyListener() *KeyListener {
	return &KeyListener{NewKeyboard(), []key.Key{}, sync.Mutex{}, nil}
}

func (k *KeyListener) Events(ctx context.Context) (chan KeyEvent, error) {
	if k.eventChan == nil {
		k.eventChan = make(chan KeyEvent)

		keyboardEventChan, err := k.keyboard.Events(ctx)
		if err != nil {
			return nil, err
		}

		go func() {
			// read events from the keyboard
			for {
				select {
				case <-ctx.Done():
					close(k.eventChan)
					k.eventChan = nil
					return
				default:
					event, ok := <-keyboardEventChan
					if !ok {
						close(k.eventChan)
						k.eventChan = nil
						return
					}

					go func(e KeyboardEvent) {
						k.currentModifiersLock.Lock()
						defer k.currentModifiersLock.Unlock()

						switch e.EventType {
						case KeyboardEventTypeDown:
							if key.IsModifierKey(e.Key) {
								k.currentModifiers = lo.Uniq(append(k.currentModifiers, e.Key))
							}
						case KeyboardEventTypeUp:
							if key.IsModifierKey(e.Key) {
								k.currentModifiers = lo.Filter(k.currentModifiers, func(code key.Key, _ int) bool {
									return code != e.Key
								})
							}

							k.eventChan <- KeyEvent{
								Key:       e.Key,
								Modifiers: k.currentModifiers,
							}
						}
					}(event)
				}
			}
		}()
	}

	return k.eventChan, nil
}
