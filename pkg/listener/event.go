package listener

import "github.com/nathan-fiscaletti/key-logger/pkg/key"

// KeyEvent represents a single key up event.
type KeyEvent struct {
	// The modifiers that were pressed at the time of the event.
	Modifiers []key.Key
	// The key that was released.
	Key key.Key
}

// KeyEventHandler is a function that handles KeyEvents.
type KeyEventHandler func(KeyEvent)
