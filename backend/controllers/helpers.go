package controllers

import "errors"

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
