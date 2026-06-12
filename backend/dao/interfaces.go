package db

import "backend/domain"

// TxContext is a marker interface for a database transaction.
// The concrete implementation (e.g. *gorm.DB) is supplied by the DAO layer.
type TxContext interface {
}

// ---------------------------------------------------------------------------
// EventDAO
// ---------------------------------------------------------------------------

type EventDAO interface {
	FindAll(filters domain.EventFilters) ([]domain.Event, error)
	FindByID(id uint) (*domain.Event, error)
	Create(event *domain.Event) error
	Update(event *domain.Event) error
	Delete(id uint) error
	IncrementTicketsSold(eventID uint, delta int) error
	DecrementTicketsSold(eventID uint) error
}

// ---------------------------------------------------------------------------
// TicketDAO
// ---------------------------------------------------------------------------

type TicketDAO interface {
	Create(ticket *domain.Ticket) error
	FindByID(id uint) (*domain.Ticket, error)
	FindByUserID(userID uint) ([]domain.Ticket, error)
	FindActiveByEvent(eventID uint) ([]domain.Ticket, error)
	CountActiveByEvent(eventID uint) (int, error)
	CancelByEvent(eventID uint) error
	Save(ticket *domain.Ticket) error
	WithTransaction(fn func(TxContext) error) error
}

// ---------------------------------------------------------------------------
// UserDAO
// ---------------------------------------------------------------------------

type UserDAO interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
}
