package errorsx

import "errors"

func Is(err, target error) bool { return errors.Is(err, target) }

// As will panic if target is not a non-nil pointer to either a type that implements
// error, or to any interface type. As returns false if err is nil.
func As(err error, target interface{}) bool { return errors.As(err, target) }

func Unwrap(err error) error {
	return errors.Unwrap(err)
}
