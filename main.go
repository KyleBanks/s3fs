// Package main acts as the primary entry-point for the gos3 application.
package main

import (
	"fmt"

	"github.com/KyleBanks/s3fs/client"
	"github.com/KyleBanks/s3fs/handler"
	"github.com/KyleBanks/s3fs/indicator"
	"github.com/KyleBanks/s3fs/listener"
)

func main() {
	// Determine the UI indicator to use.
	ui := indicator.NewCommandLine()

	// Determine the required handler and listener types.
	// Note: In the future there may be more than one kind to choose from, especially likely for the listener (ie. http listener?).
	var h handler.Handler
	var l listener.Listener

	h = handler.NewS3(client.New("us-east-1"), ui)
	l = listener.NewStdin(ui)

	// Infinitely listen for and handle input from the user.
	for {
		if cmd, ok := l.Listen(); ok {
			if err := h.Handle(cmd); err != nil {
				fmt.Println(err)
			}
		}
	}
}
