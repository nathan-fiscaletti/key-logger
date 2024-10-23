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
