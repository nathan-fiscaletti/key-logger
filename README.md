# key-logger

This is a simple implementation of a cross-platform key logger in Go.

## Supported Platforms

- Windows
- Linux (code is untested)

## Usage

### Raw Keyboard Events

Raw keyboard events will include both the key pressed and the key released events.

```go
package main

import (
	"fmt"

	"github.com/nathan-fiscaletti/key-logger/pkg/keyboard"
	"github.com/nathan-fiscaletti/key-logger/pkg/keyboard/key"

	"context"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	kb := keyboard.NewKeyboard()

	keyChannel, err := kb.Events(ctx)
	if err != nil {
		fmt.Printf("Error getting keyboard events: %v\n", err)
		return
	}

	fmt.Println("Listening for keyboard events...")
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-keyChannel:
			go func() {
				fmt.Printf("Key: %s, Event: %s\n", event.Key.Name, event.EventType)

				if event.Key.Equals(key.Escape) {
					cancel()
				}
			}()
		}
	}
}
```

### Key Listener

The key listener will only include the key released events, but each event will also include a list
of modifier keys that were pressed at the time of the event.

```go
package main

import (
	"fmt"

	"github.com/nathan-fiscaletti/key-logger/pkg/keyboard"
	"github.com/nathan-fiscaletti/key-logger/pkg/keyboard/key"
	"github.com/samber/lo"

	"context"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	keyListener := keyboard.NewKeyListener()

	keyChannel, err := keyListener.Events(ctx)
	if err != nil {
		fmt.Printf("Error getting keyboard events: %v\n", err)
		return
	}

	fmt.Println("Listening for keyboard events...")
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-keyChannel:
			go func() {
				fmt.Printf(
					"Key: %v, Modifiers: %v\n",
					event.Key.Name,
					lo.Map(event.Modifiers, func(k key.Key, _ int) string { return k.Name }),
				)
				if event.Key.Equals(key.Escape) {
					cancel()
				}
			}()
		}
	}
}
```