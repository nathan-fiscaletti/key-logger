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
