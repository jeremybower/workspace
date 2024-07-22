package internal

import (
	"errors"
	"fmt"
)

var ErrAbsolutePathRequired = fmt.Errorf("absolute path required")

var ErrConfigInvalid = fmt.Errorf("invalid config")

var ErrMountInvalid = errors.New("invalid mount")
