package errors

import (
	"errors"
	"strings"
)

// Combine takes a list of errors and combines them into a single error.
func Combine(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	errorMessages := make([]string, len(errs))
	for i, err := range errs {
		errorMessages[i] = err.Error()
	}

	return errors.New(strings.Join(errorMessages, "; "))
}
