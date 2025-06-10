// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package errorsx

import "errors"

func Is(err, target error) bool { return errors.Is(err, target) }

// As will panic if target is not a non-nil pointer to either a type that implements
// error, or to any interface type. As returns false if err is nil.
func As(err error, target interface{}) bool { return errors.As(err, target) }

func Unwrap(err error) error {
	return errors.Unwrap(err)
}
