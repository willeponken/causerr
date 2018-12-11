// Package causerr implements errors with messages and stack traces.
package causerr

import (
	"fmt"

	stderr "errors"

	"github.com/pkg/errors"
)

// causeError represents an error and a message.
type causeError struct {
	err     error
	Message string `json:"message"`
}

// Error fullfills the error interface for printing error messages together with
// messages from the causeError type.
func (d *causeError) Error() string {
	return fmt.Sprintf("%v (%s)", d.err.Error(), d.Message)
}

// Format fullfills the fmt.Formatter interface for pretty-printing the
// causeError type.
func (d *causeError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%s\n%+v", d.Message, d.err)
			return
		}
		fallthrough
	case 's':
		fmt.Fprintf(s, "%s\n%v", d.Message, d.err)
	case 'q':
		fmt.Fprintf(s, "%s\n%q", d.Message, d.err)
	}
}

// New creates an error with cause, message and stack trace.
//
// The cause can be either an error or a string which will be used as
// the internal error. If the cause is not any of the supported types it
// will panic.
//
// The message is the external error that can be shown to non-developers.
func New(cause interface{}, message string) error {
	var err error

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
