package services

import "errors"

// notFoundErr implements both the error interface and the
// notFoundError interface checked by controllers via errors.As.
type notFoundErr struct{}

func (e *notFoundErr) Error() string { return "resource not found" }

// NotFound is consumed by controllers/helpers.go isNotFound().
func (e *notFoundErr) NotFound() bool { return true }

var (
	ErrNotFound          = &notFoundErr{}
	ErrUnauthorized      = errors.New("action not authorized for this user")
	ErrEventCancelled    = errors.New("event has been cancelled")
	ErrNoCapacity        = errors.New("event has no remaining capacity")
	ErrTicketNotOwned    = errors.New("ticket does not belong to this user")
	ErrAlreadyCancelled  = errors.New("ticket is already cancelled")
	ErrInvalidTransfer      = errors.New("cannot transfer to the same user")
	ErrAlreadyTransferred  = errors.New("ticket has already been transferred")
	ErrInvalidInput      = errors.New("invalid input")
	ErrSalesNotOpen      = errors.New("sales are not yet open for this event")
	ErrPresaleCodeRequired = errors.New("presale code is required")
	ErrInvalidPresaleCode = errors.New("invalid presale code")
)
