// Package indicator provides UI indications to the user.
package indicator

// stringWriter defines an interface that can receive strings.
type stringWriter interface {
	Write(string)
}
