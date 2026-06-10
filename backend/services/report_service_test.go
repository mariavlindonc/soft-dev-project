package services

import (
	"errors"
	"testing"

	"backend/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEventReport(t *testing.T) {
	t.Run("existing event returns report", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		ticketDAO := new(MockTicketDAO)
		userDAO := new(MockUserDAO)
		svc := NewReportService(eventDAO, ticketDAO, userDAO)

		eventDAO.On("FindByID", uint(1)).Return(&domain.Event{
			ID: 1, Title: "Concert", Capacity: 100,
		}, nil)
		ticketDAO.On("CountActiveByEvent", uint(1)).Return(75, nil)
		ticketDAO.On("FindActiveByEvent", uint(1)).Return([]domain.Ticket{
			{UserID: 10}, {UserID: 20},
		}, nil)
		userDAO.On("FindByID", uint(10)).Return(&domain.User{ID: 10, Name: "Alice", Email: "a@t.com"}, nil)
		userDAO.On("FindByID", uint(20)).Return(&domain.User{ID: 20, Name: "Bob", Email: "b@t.com"}, nil)

		report, err := svc.GetEventReport(1)
		require.NoError(t, err)
		assert.Equal(t, uint(1), report.EventID)
		assert.Equal(t, "Concert", report.EventTitle)
		assert.Equal(t, 100, report.TotalCapacity)
		assert.Equal(t, 75, report.TicketsSold)
		assert.Equal(t, 75.0, report.Occupancy)
		assert.Len(t, report.Buyers, 2)
	})

	t.Run("not found returns ErrNotFound", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewReportService(eventDAO, new(MockTicketDAO), new(MockUserDAO))

		eventDAO.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

		_, err := svc.GetEventReport(99)
		assert.ErrorIs(t, err, ErrNotFound)
	})
}

func TestGetGlobalReport(t *testing.T) {
	eventDAO := new(MockEventDAO)
	ticketDAO := new(MockTicketDAO)
	svc := NewReportService(eventDAO, ticketDAO, new(MockUserDAO))

	eventDAO.On("FindAll", domain.EventFilters{}).Return([]domain.Event{
		{ID: 1, Title: "E1", Capacity: 100},
		{ID: 2, Title: "E2", Capacity: 200},
	}, nil)
	ticketDAO.On("CountActiveByEvent", uint(1)).Return(50, nil)
	ticketDAO.On("CountActiveByEvent", uint(2)).Return(100, nil)

	report, err := svc.GetGlobalReport()
	require.NoError(t, err)
	assert.Equal(t, 2, report.TotalEvents)
	assert.Equal(t, 150, report.TotalTicketsSold)
	assert.Len(t, report.EventReports, 2)
}
