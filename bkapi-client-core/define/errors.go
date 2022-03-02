package define

import (
	"errors"

	pkgErrors "github.com/pkg/errors"
)

//
var (
	// ErrTypeNotMatch defines the error which indicates the type of the value is not match.
	ErrTypeNotMatch = errors.New("type not match")
)

// ErrorWrapf annotates err with the format specifier and arguments.
var ErrorWrapf = pkgErrors.WithMessagef

// EnableStackTraceErrorWrapf enables stack trace for ErrorWrapf.
func EnableStackTraceErrorWrapf() {
	SetErrorWrapf(pkgErrors.Wrapf)
}

// SetErrorWrapf sets the ErrorWrapf.
// This function is not thread safe, do not call it in parallel.
func SetErrorWrapf(f func(err error, format string, args ...interface{}) error) {
	ErrorWrapf = f
}
