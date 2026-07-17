package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrNotFound_NotFound(t *testing.T) {
	nf := &notFoundErr{}
	assert.True(t, nf.NotFound())
}

func TestErrNotFound_Error(t *testing.T) {
	nf := &notFoundErr{}
	assert.Equal(t, "resource not found", nf.Error())
}

func TestNotFoundAsSentinel(t *testing.T) {
	assert.True(t, errors.Is(ErrNotFound, ErrNotFound))
	assert.Equal(t, "resource not found", ErrNotFound.Error())
}

func TestErrUnauthorized(t *testing.T) {
	assert.Equal(t, "action not authorized for this user", ErrUnauthorized.Error())
}

func TestErrEventCancelled(t *testing.T) {
	assert.Equal(t, "event has been cancelled", ErrEventCancelled.Error())
}

func TestErrNoCapacity(t *testing.T) {
	assert.Equal(t, "event has no remaining capacity", ErrNoCapacity.Error())
}

func TestErrTicketNotOwned(t *testing.T) {
	assert.Equal(t, "ticket does not belong to this user", ErrTicketNotOwned.Error())
}

func TestErrAlreadyCancelled(t *testing.T) {
	assert.Equal(t, "ticket is already cancelled", ErrAlreadyCancelled.Error())
}

func TestErrInvalidTransfer(t *testing.T) {
	assert.Equal(t, "cannot transfer to the same user", ErrInvalidTransfer.Error())
}

func TestErrAlreadyTransferred(t *testing.T) {
	assert.Equal(t, "ticket has already been transferred", ErrAlreadyTransferred.Error())
}

func TestErrInvalidInput(t *testing.T) {
	assert.Equal(t, "invalid input", ErrInvalidInput.Error())
}

func TestErrSalesNotOpen(t *testing.T) {
	assert.Equal(t, "sales are not yet open for this event", ErrSalesNotOpen.Error())
}

func TestErrPresaleCodeRequired(t *testing.T) {
	assert.Equal(t, "presale code is required", ErrPresaleCodeRequired.Error())
}

func TestErrInvalidPresaleCode(t *testing.T) {
	assert.Equal(t, "invalid presale code", ErrInvalidPresaleCode.Error())
}

// Verify ErrNotFound implements the notFoundError interface
func TestErrNotFoundImplementsNotFoundInterface(t *testing.T) {
	type notFounder interface {
		NotFound() bool
	}
	var nfe notFounder
	var err error = ErrNotFound
	nfe, ok := err.(notFounder)
	assert.True(t, ok, "ErrNotFound should implement notFoundError")
	assert.True(t, nfe.NotFound())
}
