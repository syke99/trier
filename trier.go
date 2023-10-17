package trier

import (
	"errors"
	"time"
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

// TryIfErr is like Try, but if an error occurs, passes it to errFn before returning
func (t *Trier) TryIfErr(errFn func(err error) error, fn func(args ...any) error, args ...any) *Trier {
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

	if t.err != nil {
		err := errFn(*t.err)
		t.err = &err
	}
	return t
}

// TryRetry is a fault-tolerant version of Try.
// If fn returns an error, it will retry to run
// fn up to limit times. If limit is less than or
// equal to zero, TryRetry will continually retry
// running fn until it doesn't error
func (t *Trier) TryRetry(limit int, fn func(args ...any) error, args ...any) *Trier {
	if t.err != nil {
		return t
	}

	switch limit <= 0 {
	case true:
		for {
			err := fn(args...)
			if err == nil {
				break
			}
		}
	case false:
		for i := 0; i < limit; i++ {
			err := fn(args...)
			if err == nil {
				break
			}

			if t.err != nil {
				*t.err = errors.Join(*t.err, err)
			} else {
				*t.err = err
			}
		}
	}

	return t
}

// TryRetryIfErr is just a combination
// of TryIfErr and TryRetry, where if
// on each iteration of retrying, if
// an error is returned, it will first
// be passes to errFn before being joined
// with previous errors
func (t *Trier) TryRetryIfErr(limit int, errFn func(err error) error, fn func(args ...any) error, args ...any) *Trier {
	if t.err != nil {
		return t
	}

	switch limit <= 0 {
	case true:
		for {
			err := fn(args...)
			if err == nil {
				break
			}
		}
	case false:
		for i := 0; i < limit; i++ {
			err := fn(args...)
			if err == nil {
				break
			}

			if t.err != nil {
				*t.err = errors.Join(*t.err, errFn(err))
			} else {
				*t.err = err
			}
		}
	}

	return t
}

// TryRetryBackoff is similar to TryRetry,
// except if limit is less than or equal
// to zero, it will create a new error with
// the value "retry backoff attempted with
// limit less than or equal to zero" and
// immediately return. Otherwise, it will
// run just like TryRetry with the added
// step of waiting for the time.Duration
// returned by the provided backoff func
// before retrying on an error
func (t *Trier) TryRetryBackoff(limit int, backoff func(i int) time.Duration, fn func(args ...any) error, args ...any) *Trier {
	if t.err != nil {
		return t
	}

	switch limit <= 0 {
	case true:
		*t.err = errors.New("retry backoff attempted with limit less than or equal to zero")
	case false:
		for i := 0; i < limit; i++ {
			err := fn(args...)
			if err == nil {
				break
			}

			if t.err != nil {
				*t.err = errors.Join(*t.err, err)
			} else {
				*t.err = err
			}

			time.Sleep(backoff(i))
		}
	}

	return t
}

// TryRetryBackoffIfErr is just a combination
// of TryIfErr and TryRetryBackoff, where if
// on each iteration of retrying, if an error
// is returned, it will first be passes to
// errFn before being joined with any previous errors
func (t *Trier) TryRetryBackoffIfErr(limit int, errFn func(err error) error, backoff func(i int) time.Duration, fn func(args ...any) error, args ...any) *Trier {
	if t.err != nil {
		return t
	}

	switch limit <= 0 {
	case true:
		*t.err = errors.New("retry backoff attempted with limit less than or equal to zero")
	case false:
		for i := 0; i < limit; i++ {
			err := fn(args...)
			if err == nil {
				break
			}

			if t.err != nil {
				*t.err = errors.Join(*t.err, errFn(err))
			} else {
				*t.err = err
			}

			time.Sleep(backoff(i))
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
