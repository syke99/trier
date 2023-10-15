package trier

import (
	"errors"
)

func NewTrier() *Trier {
	return &Trier{}
}

// Trier internally keeps track of errors
// and allows you to chain function calls
// without having to keep track of whether
// an error value is nil or not
type Trier struct {
	err *error
}

// Try checks for an existing error and if
// none exists, calls fn with the given args
// and returns the *Trier that called it. Once
// all functions have been tried, call Err() or
// UnwrapErr() to obtain any returned error(s).
// If you want to try a function, but an error
// may exist, and you want to collect multiple
// errors, use TryWrap() instead
func (t *Trier) Try(fn func(args ...any) error, args ...any) *Trier {
	if t.err != nil {
		return t
	}

	err := fn(args...)

	if err != nil {
		if t.err == nil {
			t.err = &err
		} else {
			*t.err = err
		}
	}

	return t
}

// TryJoin calls fn with the given args and
// if a previous error exists and fn returns
// an error, it will join these two errors
// together with errors.Join() to allow for
// multiple errors to be collected
func (t *Trier) TryJoin(fn func(args ...any) error, args ...any) *Trier {
	err := fn(args...)

	if t.err != nil {
		x := errors.Join(*t.err, err)
		t.err = &x
	} else {
		t.err = &err
	}

	return t
}

// IfErr calls fn if the error stored in the Trier
// is not nil. This makes wrapping errors with custom
// errors a breeze
func (t *Trier) IfErr(fn func(err error) error) *Trier {
	if t.err != nil {
		err := fn(*t.err)
		t.err = &err
	}
	return t
}

// Nil allows you to nil out an error. This way a
// single trier can be used across a codebase as
// long as you know when you are nilling out errors
func (t *Trier) Nil() *Trier {
	if t.err != nil {
		t.err = nil
	}
	return t
}

// Err returns the first error experienced,
// or any wrapped errors
func (t *Trier) Err() error {
	return *t.err
}
