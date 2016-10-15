// Package listener provides input listening functionality.
package listener

// Listener defines an interface that can listen for incoming commands.
type Listener interface {
	Listen() (cmd []string, ok bool)
}
