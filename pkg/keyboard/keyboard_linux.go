package keyboard

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/nathan-fiscaletti/key-logger/pkg/key"
)

const eventSize = 24 // Size of each event struct in bytes

type inputEvent struct {
	Time  [2]uint64 // Timestamp (seconds and microseconds)
	Type  uint16    // Event type
	Code  uint16    // Event code
	Value int32     // Event value (1 for key down, 0 for key up)
}

func init() {
	keyboard = keyboardLinux{}
}

type keyboardLinux struct {
	eventChan chan Event
}

func (k keyboardLinux) Events(ctx context.Context) (chan Event, error) {
	if k.eventChan == nil {
		k.eventChan = make(chan Event)

		device, err := findKeyboardDevice()
		if err != nil {
			return nil, err
		}
		file, err := os.Open(device)
		if err != nil {
			return nil, err
		}

		go func() {
			defer file.Close()

			for {
				event := inputEvent{}
				buffer := make([]byte, eventSize)
				_, err := file.Read(buffer)
				if err != nil {
					close(k.eventChan)
					k.eventChan = nil
					return
				}

				err = binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, &event)
				if err != nil {
					close(k.eventChan)
					k.eventChan = nil
					return
				}

				if event.Type == 1 {
					eventType := KeyboardEventTypeUp
					if event.Value == 1 {
						eventType = KeyboardEventTypeDown
					}
					k.eventChan <- Event{
						Key:       key.FindKeyCode(uint32(event.Code)),
						EventType: eventType,
						Timestamp: time.Unix(int64(event.Time[0]), int64(event.Time[1])),
					}
				}
			}
		}()
	}

	return k.eventChan, nil
}

func findKeyboardDevice() (string, error) {
	for i := 0; i < 255; i++ {
		f, err := os.Open(fmt.Sprintf("/sys/class/input/event%d/device/name", i))
		if err != nil {
			return "", err
		}

		var data []byte
		data, err = io.ReadAll(f)
		if err != nil {
			return "", err
		}
		content := string(data)

		if strings.Contains(strings.ToLower(content), "mouse") {
			continue
		}

		for _, identifier := range []string{"keyboard", "mx keys"} {
			if strings.Contains(strings.ToLower(content), identifier) {
				return fmt.Sprintf("/dev/input/event%d", i), nil
			}
		}
	}

	return "", nil
}
