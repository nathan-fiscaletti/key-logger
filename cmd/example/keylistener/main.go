package main

import (
	"fmt"
	"log"

	"github.com/nathan-fiscaletti/key-logger/pkg/key"
	"github.com/nathan-fiscaletti/key-logger/pkg/keyboard"
	"github.com/nathan-fiscaletti/key-logger/pkg/listener"
	"github.com/samber/lo"

	"context"
)

func main() {
	// Initialize a context.
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Retrieve the main keyboard hook.
	kb := keyboard.New()

	// Create a new listener using the keyboard hook.
	kl := listener.New(kb)

	// Start listening for events.
	err := kl.Listen(ctx, func(event listener.KeyEvent) {
		// Determine which modifier keys were held down.
		//
		// Technically, event.Modifiers will contain any other key that was held
		// down at the time of that event. For display purposes, we should only
		// display the names of the ones that are actually considered modifier keys.
		//
		// You can customize these keys using `key.AddModifierKey()`
		// and `key.RemoveModifierKey()`.
		//
		// You can check if a key is a modifier key using `key.IsModifierKey()`.
		//
		// By default, they are Shift, Control, Alt or Win. (Either side.)
		modifierKeys := lo.Filter(event.Modifiers, func(k key.Key, _ int) bool {
			return key.IsModifierKey(k)
		})

		// Print out a message about the key.
		fmt.Printf(
			"Key: %v, Modifiers: %v\n",

			// Print the name of the key.
			event.Key.Name,

			// Print the names of each key that is a modifier.
			lo.Map(modifierKeys, func(k key.Key, _ int) string { return k.Name }),
		)

		// If the escape key is pressed, cancel the context to end the program.
		if event.Key.Equals(key.Escape) {
			cancel()
		}
	})

	if err != nil {
		log.Fatal("Error listening for keyboard events\n", err)
		return
	}

	fmt.Println("Listening for keyboard events...")
	<-ctx.Done()
}
