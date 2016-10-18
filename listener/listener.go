// Package listener provides input listening functionality.
package listener

// Listener defines an interface that can listen for incoming commands.
type Listener interface {
	Listen() (cmds []InputCommand, ok bool)
}

// InputCommand defines a command input recieved by the listener.
type InputCommand struct {
	Args []string
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
