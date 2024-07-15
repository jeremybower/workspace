package funcs

import (
	"errors"
	"reflect"
)

var ErrValueRequired = errors.New("zero value is not allowed")

func Require() func(value any) (any, error) {
	return func(value any) (any, error) {
		if value == nil {
			return nil, ErrValueRequired
		}

		v := reflect.ValueOf(value)
		zero := reflect.Zero(v.Type())
		if reflect.DeepEqual(value, zero.Interface()) {
			return nil, ErrValueRequired
		}

		return value, nil
	}
}
