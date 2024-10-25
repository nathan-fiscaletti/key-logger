# key-logger

[![Sponsor Me!](https://img.shields.io/badge/%F0%9F%92%B8-Sponsor%20Me!-blue)](https://github.com/sponsors/nathan-fiscaletti)
[![GoDoc](https://godoc.org/github.com/nathan-fiscaletti/key-logger?status.svg)](https://godoc.org/github.com/nathan-fiscaletti/key-logger)

This is a simple implementation of a cross-platform key logger in Go.

## Supported Platforms

- Windows
- Linux

## Usage

```sh
go get github.com/nathan-fiscaletti/key-logger
```

### Raw Keyboard Events

Raw keyboard events will include both the key pressed and the key released events.

Run this example with: `go run ./cmd/example/raw`
```go
package main

import (
	"fmt"

	"github.com/nathan-fiscaletti/key-logger/pkg/key"
	"github.com/nathan-fiscaletti/key-logger/pkg/keyboard"

	"context"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	kb := keyboard.New()

	err := kb.Listen(ctx, func(event keyboard.Event) {
		fmt.Printf("Key: %s, Event: %s\n", event.Key.Name, event.EventType)

		// If the escape key is pressed, cancel the context
		if event.Key.Equals(key.Escape) {
			cancel()
		}
	})

	if err != nil {
		fmt.Printf("Error listening for keyboard events: %v\n", err)
		return
	}

	fmt.Println("Listening for keyboard events...")
	<-ctx.Done()
}
```

### Key Listener

The key listener will only include the key released events, but each event will also include a list
of modifier keys that were pressed at the time of the event.

You can use `key.IsModifierKey` to check if a key is a modifier key.

Run this example with: `go run ./cmd/example/keylistener`
```go
package main

import (
	"fmt"

	"github.com/nathan-fiscaletti/key-logger/pkg/key"
	"github.com/nathan-fiscaletti/key-logger/pkg/keyboard"
	"github.com/nathan-fiscaletti/key-logger/pkg/listener"
	"github.com/samber/lo"

	"context"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	kb := keyboard.New()
	kl := listener.New(kb)

	err := kl.Listen(ctx, func(event listener.KeyEvent) {
		fmt.Printf(
			"Key: %v, Modifiers: %v\n",
			event.Key.Name,
			lo.Map(
				lo.Filter(event.Modifiers, func(k key.Key, _ int) bool {
					return key.IsModifierKey(k)
				}),
				func(k key.Key, _ int) string { return k.Name },
			),
		)

		// If the escape key is pressed, cancel the context
		if event.Key.Equals(key.Escape) {
			cancel()
		}
	})

	if err != nil {
		fmt.Printf("Error listening for keyboard events: %v\n", err)
		return
	}

	fmt.Println("Listening for keyboard events...")
	<-ctx.Done()
}
```