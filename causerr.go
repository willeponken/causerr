// Package causerr implements errors with IDs, messages and stack trace.
package causerr

import (
	"fmt"

	stderr "errors"

	"github.com/pkg/errors"
)

// causeError represents an error with ID and a message.
type causeError struct {
	err     error
	ID      int    `json:"id"`
	Message string `json:"message"`
}

// Error fullfills the error interface for printing error messages together with ID and
// messages from the causeError type.
func (d *causeError) Error() string {
	return fmt.Sprintf("%v (%d: %s)", d.err.Error(), d.ID, d.Message)
}

// Format fullfills the fmt.Formatter interface for pretty-printing the causeError type.
func (d *causeError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "#%d: %s\n%+v", d.ID, d.Message, d.err)
			return
		}
		fallthrough
	case 's':
		fmt.Fprintf(s, "#%d: %s\n%v", d.ID, d.Message, d.err)
	case 'q':
		fmt.Fprintf(s, "#%d: %s\n%q", d.ID, d.Message, d.err)
	}
}

// New creates an error with ID, cause, message and stack trace.
//
// The cause can be either an error or a string which will be used as
// the internal error. If the cause is not any of the supported types it
// will panic.
//
// The ID is an internal number to identify errors, it must be >=0, else
// it will panic.
//
// The message is the external error that can be shown to non-developers.
func New(id int, cause interface{}, message string) error {
	var err error

	if id < 0 {
		panic("id must be >=0")
	}

	switch cause.(type) {
	case error:
		err = cause.(error)
	case string:
		err = stderr.New(cause.(string))
	default:
		panic(fmt.Sprintf("invalid type for cause: %#v (%T)", cause, cause))
	}

	return &causeError{
		err:     errors.WithStack(err),
		ID:      id,
		Message: message,
	}
}

// getCauseError is used internally to type assert into the causeError type.
func getCauseError(err error) (*causeError, bool) {
	def, ok := err.(*causeError)
	if ok {
		return def, true
	}

	return nil, false
}

// ID takes an error and returns its ID.
// If no ID could be found it'll return -1.
func ID(err error) int {
	def, ok := getCauseError(err)
	if ok {
		return def.ID
	}

	return -1
}

// Cause takes an error and returns its error cause.
// If no cause can be found it'll return nil.
func Cause(err error) error {
	def, ok := getCauseError(err)
	if ok {
		return def.err
	}

	return nil
}

// Message takes an error and returns its message.
// If no message can be found it'll return an empty string.
func Message(err error) string {
	def, ok := getCauseError(err)
	if ok {
		return def.Message
	}

	return ""
}
