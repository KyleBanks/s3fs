// Package handler provides command processing functionality.
package handler

import (
	"github.com/KyleBanks/s3fs/handler/command"
)

// Handler defines an interface that can accept and process commands.
type Handler interface {
	Handle(cmd []string, out command.Outputter) error
}
