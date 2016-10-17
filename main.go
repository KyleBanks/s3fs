// Package main acts as the primary entry-point for the gos3 application.
package main

import (
	"bufio"
	"os"

	"github.com/KyleBanks/s3fs/client"
	"github.com/KyleBanks/s3fs/handler"
	"github.com/KyleBanks/s3fs/indicator"
	"github.com/KyleBanks/s3fs/listener"
	"github.com/KyleBanks/s3fs/output"
)

func main() {
	// Determine the output method to use.
	out := output.New(os.Stdout)

	// Determine the UI indicator to use.
	ui := indicator.NewCommandLine(out)

	// Determine the required handler and listener types.
	// Note: In the future there may be more than one kind to choose from, especially likely for the listener (ie. http listener?).
	var h handler.Handler
	var l listener.Listener

	h = handler.NewS3(client.New("us-east-1"), ui)
	l = listener.NewText(ui, bufio.NewScanner(os.Stdin))

	// Infinitely listen for and handle input from the user.
	for {
		if cmd, ok := l.Listen(); ok {
			if err := h.Handle(cmd, out); err != nil {
				out.Write(err.Error() + "\n")
			}
		}
	}
}
