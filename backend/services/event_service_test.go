package services

import (
	"errors"
	"testing"
	"time"

	"backend/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestEventGetAll(t *testing.T) {
	eventDAO := new(MockEventDAO)
	svc := NewEventService(eventDAO, new(MockTicketDAO))

	expected := []domain.Event{{Title: "Concert"}, {Title: "Festival"}}
	eventDAO.On("FindAll", mock.Anything).Return(expected, nil)

	events, err := svc.GetAll(domain.EventFilters{})
	require.NoError(t, err)
	assert.Len(t, events, 2)
	eventDAO.AssertExpectations(t)
}

func TestEventGetByID(t *testing.T) {
	t.Run("found returns event", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		eventDAO.On("FindByID", uint(1)).Return(&domain.Event{Title: "Concert"}, nil)

		event, err := svc.GetByID(1)
		require.NoError(t, err)
		assert.Equal(t, "Concert", event.Title)
	})

	t.Run("not found returns ErrNotFound", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		eventDAO.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

		_, err := svc.GetByID(99)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}

func TestEventCreate(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)

	t.Run("valid event creates successfully", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		eventDAO.On("Create", mock.MatchedBy(func(e *domain.Event) bool {
			return e.Title == "New Event" && e.Capacity == 100
		})).Return(nil)

		event, err := svc.Create(CreateEventInput{
			Title:    "New Event",
			Date:     future,
			Capacity: 100,
			Price:    50.0,
		})
		require.NoError(t, err)
		assert.Equal(t, "New Event", event.Title)
		eventDAO.AssertExpectations(t)
	})

	t.Run("empty title returns ErrInvalidInput", func(t *testing.T) {
		svc := NewEventService(new(MockEventDAO), new(MockTicketDAO))
		_, err := svc.Create(CreateEventInput{
			Title:    "",
			Date:     future,
			Capacity: 100,
		})
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("capacity zero returns ErrInvalidInput", func(t *testing.T) {
		svc := NewEventService(new(MockEventDAO), new(MockTicketDAO))
		_, err := svc.Create(CreateEventInput{
			Title:    "Event",
			Date:     future,
			Capacity: 0,
		})
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("past date returns ErrInvalidInput", func(t *testing.T) {
		svc := NewEventService(new(MockEventDAO), new(MockTicketDAO))
		_, err := svc.Create(CreateEventInput{
			Title:    "Event",
			Date:     time.Now().Add(-1 * time.Hour),
			Capacity: 100,
		})
		assert.ErrorIs(t, err, ErrInvalidInput)
	})
}

func TestEventCancel(t *testing.T) {
	t.Run("cancel active event succeeds", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		ticketDAO := new(MockTicketDAO)
		svc := NewEventService(eventDAO, ticketDAO)

		event := &domain.Event{ID: 1, Status: "active"}
		eventDAO.On("FindByID", uint(1)).Return(event, nil)
		eventDAO.On("Update", mock.MatchedBy(func(e *domain.Event) bool {
			return e.Status == "cancelled"
		})).Return(nil)
		ticketDAO.On("CancelByEvent", uint(1)).Return(nil)

		err := svc.Cancel(1)
		require.NoError(t, err)
		eventDAO.AssertExpectations(t)
		ticketDAO.AssertExpectations(t)
	})

	t.Run("cancel already cancelled returns ErrEventCancelled", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		eventDAO.On("FindByID", uint(1)).Return(&domain.Event{ID: 1, Status: "cancelled"}, nil)

		err := svc.Cancel(1)
		assert.ErrorIs(t, err, ErrEventCancelled)
	})

	t.Run("cancel not found returns ErrNotFound", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		eventDAO.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

		err := svc.Cancel(99)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}

func TestEventUpdate(t *testing.T) {
	t.Run("update existing event succeeds", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		existing := &domain.Event{ID: 1, Title: "Old", Status: "active"}
		eventDAO.On("FindByID", uint(1)).Return(existing, nil)

		newTitle := "Updated"
		eventDAO.On("Update", mock.MatchedBy(func(e *domain.Event) bool {
			return e.Title == "Updated"
		})).Return(nil)

		event, err := svc.Update(1, UpdateEventInput{Title: &newTitle})
		require.NoError(t, err)
		assert.Equal(t, "Updated", event.Title)
		eventDAO.AssertExpectations(t)
	})

	t.Run("update cancelled event returns ErrEventCancelled", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		eventDAO.On("FindByID", uint(1)).Return(&domain.Event{ID: 1, Status: "cancelled"}, nil)

		_, err := svc.Update(1, UpdateEventInput{})
		assert.ErrorIs(t, err, ErrEventCancelled)
	})

	t.Run("update not found returns ErrNotFound", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		eventDAO.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

		_, err := svc.Update(99, UpdateEventInput{})
		assert.ErrorIs(t, err, ErrNotFound)
	})
}
