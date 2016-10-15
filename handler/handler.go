// Package handler provides command processing functionality.
package handler

// Handler defines an interface that can accept and process commands.
type Handler interface {
	Handle(cmd []string) error
}
