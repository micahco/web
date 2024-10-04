package httperr

import (
	"fmt"
	"net/http"
)

// Some common errors with default status text message.
var (
	ErrNotFound     = errNotFound()
	ErrUnauthorized = errUnauthorized()
)

type Error struct {
	// Wrapped error
	err error

	// Use http.StatusCode
	code int

	// Message is user facing and therefore should only
	// contain information relevant to the user.
	message string
}

func (e Error) Error() string {
	return fmt.Sprintf("httperr [%d]: %s", e.code, e.message)
}

func (e Error) Unwrap() error {
	return e.err
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Message() string {
	if e.message != "" {
		return e.message
	}

	return http.StatusText(e.code)
}

// Wrap error with relevant user facing message.
func WrapWithMessage(err error, code int, message string) error {
	return &Error{
		err:     err,
		code:    code,
		message: message,
	}
}

// Wrap error with default status text message.
func Wrap(err error, code int) error {
	return WrapWithMessage(err, code, http.StatusText(code))
}

// Create httperr without wrapping.
func New(code int, message string) error {
	return &Error{
		code:    code,
		message: message,
	}
}

// Create new error with default status text message.
func newFromStatusCode(code int) error {
	return New(code, http.StatusText(code))
}

func errNotFound() error {
	return newFromStatusCode(http.StatusNotFound)
}

func errUnauthorized() error {
	return newFromStatusCode(http.StatusUnauthorized)
}
