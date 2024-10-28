package keyboard

import (
	"context"
	"fmt"
	"time"

	input "github.com/nathan-fiscaletti/dev-input"

	"github.com/nathan-fiscaletti/key-logger/pkg/key"
)

// On linux there is a special event type for repeat events.
const KeyboardEventTypeRepeat EventType = "repeat"

// const eventSize = 24 // Size of each event struct in bytes
func init() {
	keyboard = keyboardLinux{}
}

type keyboardLinux struct {
	eventChan chan Event
}

func (k keyboardLinux) Events(ctx context.Context) (chan Event, error) {
	if k.eventChan == nil {
		k.eventChan = make(chan Event)

		keyboards, err := input.ListKeyboards()
		if err != nil {
			return nil, err
		}

		if len(keyboards) < 1 {
			return nil, fmt.Errorf("no keyboards found")
		}

		for _, keyboard := range keyboards {
			go func() {
				err := keyboard.Listen(ctx, func(e input.Event) {
					if e.Type == input.EV_TYPE_KEY {
						var eventType EventType
						switch e.Value {
						case 1:
							eventType = KeyboardEventTypeDown
						case 2:
							eventType = KeyboardEventTypeRepeat
						default:
							eventType = KeyboardEventTypeUp
						}
						k.eventChan <- Event{
							Key:       key.FindKeyCode(uint32(e.Code)),
							EventType: eventType,
							Timestamp: time.Unix(int64(e.Time[0]), int64(e.Time[1])),
						}
					}
				})
				if err != nil {
					fmt.Printf("Error listening for keyboard events: %v\n", err)
				}
			}()
		}
	}

	return k.eventChan, nil
}
