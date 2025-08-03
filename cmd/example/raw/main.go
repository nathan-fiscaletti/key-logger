package main

import (
	"fmt"

	"github.com/nathan-fiscaletti/key-logger/pkg/key"
	"github.com/nathan-fiscaletti/key-logger/pkg/keyboard"

	"context"
)

func main() {
	// Initialize a context.
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Retrieve the main keyboard hook.
	kb := keyboard.New()

	// Since we are using the main keyboard hook, we will receive both up and
	// down events for every key on the keyboard.
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
