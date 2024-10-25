package keyboard

import (
	"time"

	"github.com/nathan-fiscaletti/key-logger/pkg/key"
)

type EventType string

const (
	KeyboardEventTypeUp   EventType = "up"
	KeyboardEventTypeDown EventType = "down"
)

// Event represents a single keyboard event.
type Event struct {
	// The key that was pressed or released.
	Key key.Key
	// The type of event that occurred. (up or down)
	EventType EventType
	// The time the event occurred.
	Timestamp time.Time
}

// EventHandler is a function that handles KeyboardEvents.
type EventHandler func(Event)
