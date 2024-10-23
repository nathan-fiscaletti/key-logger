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
