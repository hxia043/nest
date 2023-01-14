// Package exception define all self exception types.
// This file defines errors which about users.
package exception

import (
	"fmt"
	"strings"
)

type InvalidImageNameError struct {
	baseError
	names []string
}

func NewInvalidImageNameError(names ...string) InvalidImageNameError {
	if len(names) < 1 {
		panic("empty names received")
	}

	message := fmt.Sprintf("invalid image name: %s", strings.Join(names, ","))
	e := InvalidImageNameError{names: names}
	e.baseError = newBaseError(message)
	return e
}
