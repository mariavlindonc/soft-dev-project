package services

import (
	"fmt"
	"time"

	db "backend/dao"
	"backend/domain"
)

// ---------------------------------------------------------------------------
// Interface
// ---------------------------------------------------------------------------

type EventServiceInterface interface {
	GetAll(filters domain.EventFilters) ([]domain.Event, error)
	GetByID(id uint) (*domain.Event, error)
	Create(input CreateEventInput) (*domain.Event, error)
	Update(id uint, input UpdateEventInput) (*domain.Event, error)
	Cancel(id uint) error
}

// ---------------------------------------------------------------------------
// Input types
// ---------------------------------------------------------------------------

type CreateEventInput struct {
	Title            string
	Date             time.Time
	Duration         int
	Capacity         int
	Price            float64
	Category         string
	Description      string
	Location         string
	ImageURL         string
	PresaleCode      string
	PresaleStartDate *time.Time
	GeneralSaleDate  *time.Time
	CreatedByID      uint
}

type UpdateEventInput struct {
	Title            *string
	Date             *time.Time
	Duration         *int
	Capacity         *int
	Price            *float64
	Category         *string
	Description      *string
	Location         *string
	ImageURL         *string
	PresaleCode      *string
	PresaleStartDate *string
	GeneralSaleDate  *string
	Status           *string
}

// ---------------------------------------------------------------------------
// Concrete service
// ---------------------------------------------------------------------------

type EventService struct {
	eventDAO  db.EventDAO
	ticketDAO db.TicketDAO
}

func NewEventService(eventDAO db.EventDAO, ticketDAO db.TicketDAO) *EventService {
	return &EventService{
		eventDAO:  eventDAO,
		ticketDAO: ticketDAO,
	}
}

func (s *EventService) GetAll(filters domain.EventFilters) ([]domain.Event, error) {
	return s.eventDAO.FindAll(filters)
}

func (s *EventService) GetByID(id uint) (*domain.Event, error) {
	event, err := s.eventDAO.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("event %d: %w", id, ErrNotFound)
	}
	return event, nil
}

func (s *EventService) Create(input CreateEventInput) (*domain.Event, error) {
	if input.Title == "" {
		return nil, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}
	if input.Capacity <= 0 {
		return nil, fmt.Errorf("%w: capacity must be greater than 0", ErrInvalidInput)
	}
	if !input.Date.After(time.Now()) {
		return nil, fmt.Errorf("%w: event date must be in the future", ErrInvalidInput)
	}

	event := &domain.Event{
		Title:            input.Title,
		Description:      strPtr(input.Description),
		ImageURL:         strPtr(input.ImageURL),
		Category:         strPtr(input.Category),
		Location:         strPtr(input.Location),
		EventDate:        input.Date,
		DurationMinutes:  input.Duration,
		Capacity:         input.Capacity,
		Price:            input.Price,
		PresaleCode:      strPtr(input.PresaleCode),
		PresaleStartDate: input.PresaleStartDate,
		GeneralSaleDate:  input.GeneralSaleDate,
		CreatedByID:      input.CreatedByID,
	}

	if err := s.eventDAO.Create(event); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return event, nil
}

func (s *EventService) Update(id uint, input UpdateEventInput) (*domain.Event, error) {
	event, err := s.eventDAO.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("event %d: %w", id, ErrNotFound)
	}

	if event.Status == "cancelled" {
		return nil, fmt.Errorf("%w: cannot update a cancelled event", ErrEventCancelled)
	}

	if input.Title != nil {
		event.Title = *input.Title
	}
	if input.Description != nil {
		event.Description = input.Description
	}
	if input.ImageURL != nil {
		event.ImageURL = input.ImageURL
	}
	if input.Category != nil {
		event.Category = input.Category
	}
	if input.Location != nil {
		event.Location = input.Location
	}
	if input.Date != nil {
		event.EventDate = *input.Date
	}
	if input.Duration != nil {
		event.DurationMinutes = *input.Duration
	}
	if input.Capacity != nil {
		event.Capacity = *input.Capacity
	}
	if input.Price != nil {
		event.Price = *input.Price
	}
	if input.PresaleCode != nil {
		event.PresaleCode = input.PresaleCode
	}
	if input.PresaleStartDate != nil {
		t, err := time.Parse(time.RFC3339, *input.PresaleStartDate)
		if err == nil {
			event.PresaleStartDate = &t
		}
	}
	if input.GeneralSaleDate != nil {
		t, err := time.Parse(time.RFC3339, *input.GeneralSaleDate)
		if err == nil {
			event.GeneralSaleDate = &t
		}
	}
	if input.Status != nil {
		event.Status = *input.Status
	}

	if err := s.eventDAO.Update(event); err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return event, nil
}

func (s *EventService) Cancel(id uint) error {
	event, err := s.eventDAO.FindByID(id)
	if err != nil {
		return fmt.Errorf("event %d: %w", id, ErrNotFound)
	}

	if event.Status == "cancelled" {
		return fmt.Errorf("event %d: %w", id, ErrEventCancelled)
	}

	event.Status = "cancelled"
	if err := s.eventDAO.Update(event); err != nil {
		return fmt.Errorf("failed to cancel event: %w", err)
	}

	// Cancel all active tickets for this event.
	if err := s.ticketDAO.CancelByEvent(event.ID); err != nil {
		return fmt.Errorf("failed to cancel tickets for event %d: %w", event.ID, err)
	}

	return nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

var _ EventServiceInterface = (*EventService)(nil)
