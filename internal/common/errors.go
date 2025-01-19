package common

import (
	"errors"
	"strings"
)

const (
	BeginTsCannotBeInTheFutureMsg = "begin timestamp cannot be in the future"
	EndTsCannotBeInTheFutureMsg   = "end timestamp cannot be in the future"
)

var (
	ErrBeginTsCannotBeInTheFuture = errors.New(BeginTsCannotBeInTheFutureMsg)
	ErrEndTsCannotBeInTheFuture   = errors.New(EndTsCannotBeInTheFutureMsg)
)

type MultiError struct {
	errors []error
}

func (m *MultiError) Error() string {
	errorMessages := make([]string, len(m.errors))
	for i, err := range m.errors {
		errorMessages[i] = "- " + err.Error()
	}
	return "\n" + strings.Join(errorMessages, "\n")
}

func (m *MultiError) Unwrap() error {
	if len(m.errors) == 0 {
		return nil
	}
	return m.errors[0]
}

func (m *MultiError) Is(target error) bool {
	for _, err := range m.errors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (m *MultiError) As(target interface{}) bool {
	for _, err := range m.errors {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

func (m *MultiError) Add(err error) {
	if err != nil {
		m.errors = append(m.errors, err)
	}
}

func (m *MultiError) NumErrors() int {
	return len(m.errors)
}
