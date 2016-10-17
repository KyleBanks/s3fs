// Package listener provides input listening functionality.
package listener

// Listener defines an interface that can listen for incoming commands.
type Listener interface {
	Listen() (cmd []string, ok bool)
}

// indicator defines a UI interface to display status updates to the user.
type indicator interface {
	ShowPrompt()
}

// inputter defines an interface that can scan for and retrieve input.
type inputter interface {
	Scan() bool
	Text() string
}
