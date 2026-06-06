package controllers

import (
	"errors"
	"time"
)

// notFoundError is implemented by service-layer error types
// to signal a resource-not-found condition without coupling
// controllers to service-specific sentinel values.
type notFoundError interface {
	NotFound() bool
}

func isNotFound(err error) bool {
	var nf notFoundError
	return errors.As(err, &nf)
}

func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func timePtrToString(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}
