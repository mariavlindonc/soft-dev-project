package services

import (
	"errors"
	"fmt"
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

	t.Run("with presale validation error propagates", func(t *testing.T) {
		svc := NewEventService(new(MockEventDAO), new(MockTicketDAO))
		_, err := svc.Create(CreateEventInput{
			Title:         "Event",
			Date:          time.Now().Add(48 * time.Hour),
			Capacity:      100,
			PresaleActive: true,
		})
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("with presale active creates successfully", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		eventDAO.On("Create", mock.AnythingOfType("*domain.Event")).Return(nil)

		presaleStart := time.Now().Add(1 * time.Hour)
		generalSale := time.Now().Add(24 * time.Hour)
		event, err := svc.Create(CreateEventInput{
			Title:            "Presale Event",
			Date:             time.Now().Add(48 * time.Hour),
			Capacity:         100,
			Price:            50,
			PresaleActive:    true,
			PresaleCode:      "CODE123",
			PresaleStartDate: &presaleStart,
			GeneralSaleDate:  &generalSale,
		})
		require.NoError(t, err)
		assert.True(t, event.PresaleActive)
		assert.Equal(t, "CODE123", *event.PresaleCode)
	})

	t.Run("dao error on create propagates", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		eventDAO.On("Create", mock.AnythingOfType("*domain.Event")).Return(fmt.Errorf("db error"))

		_, err := svc.Create(CreateEventInput{
			Title:    "Event",
			Date:     time.Now().Add(48 * time.Hour),
			Capacity: 100,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("with optional string fields creates successfully", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		eventDAO.On("Create", mock.AnythingOfType("*domain.Event")).Return(nil)

		desc := "A great show"
		event, err := svc.Create(CreateEventInput{
			Title:       "Concert",
			Date:        time.Now().Add(48 * time.Hour),
			Capacity:    100,
			Description: desc,
			Category:    "Rock",
			Location:    "Stadium",
			ImageURL:    "http://img.jpg",
		})
		require.NoError(t, err)
		assert.Equal(t, "A great show", *event.Description)
		assert.Equal(t, "Rock", *event.Category)
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

	t.Run("cancel Update error propagates", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		event := &domain.Event{ID: 1, Status: "active"}
		eventDAO.On("FindByID", uint(1)).Return(event, nil)
		eventDAO.On("Update", mock.Anything).Return(fmt.Errorf("update error"))

		err := svc.Cancel(1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update error")
	})

	t.Run("cancel ticket cancellation error propagates", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		ticketDAO := new(MockTicketDAO)
		svc := NewEventService(eventDAO, ticketDAO)

		event := &domain.Event{ID: 1, Status: "active"}
		eventDAO.On("FindByID", uint(1)).Return(event, nil)
		eventDAO.On("Update", mock.Anything).Return(nil)
		ticketDAO.On("CancelByEvent", uint(1)).Return(fmt.Errorf("db error"))

		err := svc.Cancel(1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
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

	t.Run("update multiple fields succeeds", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		existing := &domain.Event{ID: 1, Title: "Old", Capacity: 50, Price: 10.0, Status: "active"}
		eventDAO.On("FindByID", uint(1)).Return(existing, nil)

		newTitle := "New Title"
		newCat := "Rock"
		newLoc := "Buenos Aires"
		newDesc := "A great show"
		newImg := "http://img.jpg"
		cap := 200
		price := 99.9
		status := "presale"
		eventDAO.On("Update", mock.MatchedBy(func(e *domain.Event) bool {
			return e.Title == "New Title" && e.Capacity == 200 && e.Price == 99.9
		})).Return(nil)

		event, err := svc.Update(1, UpdateEventInput{
			Title:       &newTitle,
			Category:    &newCat,
			Location:    &newLoc,
			Description: &newDesc,
			ImageURL:    &newImg,
			Capacity:    &cap,
			Price:       &price,
			Status:      &status,
		})
		require.NoError(t, err)
		assert.Equal(t, "New Title", event.Title)
		assert.Equal(t, 200, event.Capacity)
	})

	t.Run("update with presale active succeeds", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		eventDate := time.Now().Add(48 * time.Hour)
		existing := &domain.Event{ID: 1, Status: "active", EventDate: eventDate}
		eventDAO.On("FindByID", uint(1)).Return(existing, nil)

		presaleActive := true
		presaleCode := "CODE123"
		presaleStart := time.Now().Add(1 * time.Hour)
		generalSale := time.Now().Add(24 * time.Hour)
		presaleStartStr := presaleStart.Format(time.RFC3339)
		generalSaleStr := generalSale.Format(time.RFC3339)

		eventDAO.On("Update", mock.AnythingOfType("*domain.Event")).Return(nil)

		event, err := svc.Update(1, UpdateEventInput{
			PresaleActive:    &presaleActive,
			PresaleCode:      &presaleCode,
			PresaleStartDate: &presaleStartStr,
			GeneralSaleDate:  &generalSaleStr,
		})
		require.NoError(t, err)
		assert.True(t, event.PresaleActive)
		assert.Equal(t, "CODE123", *event.PresaleCode)
	})

	t.Run("update presale active missing fields returns ErrInvalidInput", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		eventDate := time.Now().Add(48 * time.Hour)
		existing := &domain.Event{ID: 1, Status: "active", EventDate: eventDate}
		eventDAO.On("FindByID", uint(1)).Return(existing, nil)

		presaleActive := true
		eventDAO.On("Update", mock.Anything).Return(nil)

		_, err := svc.Update(1, UpdateEventInput{
			PresaleActive: &presaleActive,
		})
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("update presale active invalid presale_start_date returns ErrInvalidInput", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		eventDate := time.Now().Add(48 * time.Hour)
		existing := &domain.Event{ID: 1, Status: "active", EventDate: eventDate}
		eventDAO.On("FindByID", uint(1)).Return(existing, nil)

		presaleActive := true
		presaleCode := "CODE"
		badDate := "not-a-date"
		goodDate := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

		_, err := svc.Update(1, UpdateEventInput{
			PresaleActive:    &presaleActive,
			PresaleCode:      &presaleCode,
			PresaleStartDate: &badDate,
			GeneralSaleDate:  &goodDate,
		})
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("update presale active invalid general_sale_date returns ErrInvalidInput", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		eventDate := time.Now().Add(48 * time.Hour)
		existing := &domain.Event{ID: 1, Status: "active", EventDate: eventDate}
		eventDAO.On("FindByID", uint(1)).Return(existing, nil)

		presaleActive := true
		presaleCode := "CODE"
		goodStart := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
		badDate := "not-a-date"

		_, err := svc.Update(1, UpdateEventInput{
			PresaleActive:    &presaleActive,
			PresaleCode:      &presaleCode,
			PresaleStartDate: &goodStart,
			GeneralSaleDate:  &badDate,
		})
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("update presale start after general sale returns ErrInvalidInput", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		eventDate := time.Now().Add(48 * time.Hour)
		existing := &domain.Event{ID: 1, Status: "active", EventDate: eventDate}
		eventDAO.On("FindByID", uint(1)).Return(existing, nil)

		presaleActive := true
		presaleCode := "CODE"
		presaleStart := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
		generalSale := time.Now().Add(1 * time.Hour).Format(time.RFC3339)

		_, err := svc.Update(1, UpdateEventInput{
			PresaleActive:    &presaleActive,
			PresaleCode:      &presaleCode,
			PresaleStartDate: &presaleStart,
			GeneralSaleDate:  &generalSale,
		})
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("update presale general sale after event date returns ErrInvalidInput", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		eventDate := time.Now().Add(2 * time.Hour)
		existing := &domain.Event{ID: 1, Status: "active", EventDate: eventDate}
		eventDAO.On("FindByID", uint(1)).Return(existing, nil)

		presaleActive := true
		presaleCode := "CODE"
		presaleStart := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
		generalSale := time.Now().Add(3 * time.Hour).Format(time.RFC3339)

		_, err := svc.Update(1, UpdateEventInput{
			PresaleActive:    &presaleActive,
			PresaleCode:      &presaleCode,
			PresaleStartDate: &presaleStart,
			GeneralSaleDate:  &generalSale,
		})
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("update presale active false clears presale fields", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		code := "OLD"
		start := time.Now()
		general := time.Now().Add(time.Hour)
		existing := &domain.Event{
			ID: 1, Status: "active",
			PresaleActive: true, PresaleCode: &code,
			PresaleStartDate: &start, GeneralSaleDate: &general,
		}
		eventDAO.On("FindByID", uint(1)).Return(existing, nil)

		presaleActive := false
		eventDAO.On("Update", mock.MatchedBy(func(e *domain.Event) bool {
			return !e.PresaleActive && e.PresaleCode == nil && e.PresaleStartDate == nil && e.GeneralSaleDate == nil
		})).Return(nil)

		event, err := svc.Update(1, UpdateEventInput{PresaleActive: &presaleActive})
		require.NoError(t, err)
		assert.False(t, event.PresaleActive)
		assert.Nil(t, event.PresaleCode)
	})

	t.Run("update individual presale fields without toggling", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		existing := &domain.Event{ID: 1, Status: "active"}
		eventDAO.On("FindByID", uint(1)).Return(existing, nil)

		newCode := "NEWCODE"
		newStart := time.Now().Add(5 * time.Hour).Format(time.RFC3339)
		newGeneral := time.Now().Add(10 * time.Hour).Format(time.RFC3339)

		eventDAO.On("Update", mock.AnythingOfType("*domain.Event")).Return(nil)

		event, err := svc.Update(1, UpdateEventInput{
			PresaleCode:      &newCode,
			PresaleStartDate: &newStart,
			GeneralSaleDate:  &newGeneral,
		})
		require.NoError(t, err)
		assert.Equal(t, "NEWCODE", *event.PresaleCode)
	})

	t.Run("update individual presale start date with invalid format is silently ignored", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		existing := &domain.Event{ID: 1, Status: "active"}
		eventDAO.On("FindByID", uint(1)).Return(existing, nil)

		badDate := "bad-format"
		eventDAO.On("Update", mock.AnythingOfType("*domain.Event")).Return(nil)

		event, err := svc.Update(1, UpdateEventInput{
			PresaleStartDate: &badDate,
		})
		require.NoError(t, err)
		assert.Nil(t, event.PresaleStartDate)
	})

	t.Run("update date and duration succeeds", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		existing := &domain.Event{ID: 1, Status: "active"}
		eventDAO.On("FindByID", uint(1)).Return(existing, nil)

		newDate := time.Now().Add(72 * time.Hour)
		newDuration := 180
		eventDAO.On("Update", mock.AnythingOfType("*domain.Event")).Return(nil)

		event, err := svc.Update(1, UpdateEventInput{
			Date:     &newDate,
			Duration: &newDuration,
		})
		require.NoError(t, err)
		assert.True(t, event.EventDate.Equal(newDate))
		assert.Equal(t, 180, event.DurationMinutes)
	})

	t.Run("update dao error returns error", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		existing := &domain.Event{ID: 1, Status: "active"}
		eventDAO.On("FindByID", uint(1)).Return(existing, nil)

		newTitle := "Updated"
		eventDAO.On("Update", mock.Anything).Return(fmt.Errorf("db error"))

		_, err := svc.Update(1, UpdateEventInput{Title: &newTitle})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("update general sale date with invalid format is silently ignored", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewEventService(eventDAO, new(MockTicketDAO))

		existing := &domain.Event{ID: 1, Status: "active"}
		eventDAO.On("FindByID", uint(1)).Return(existing, nil)

		badDate := "bad-format"
		eventDAO.On("Update", mock.AnythingOfType("*domain.Event")).Return(nil)

		event, err := svc.Update(1, UpdateEventInput{
			GeneralSaleDate: &badDate,
		})
		require.NoError(t, err)
		assert.Nil(t, event.GeneralSaleDate)
	})
}

func TestValidatePresaleConfig(t *testing.T) {
	eventDate := time.Now().Add(48 * time.Hour)

	t.Run("presale not active returns nil", func(t *testing.T) {
		err := validatePresaleConfig(false, "", nil, nil, &eventDate)
		assert.NoError(t, err)
	})

	t.Run("presale active with nil start date returns ErrInvalidInput", func(t *testing.T) {
		general := time.Now().Add(24 * time.Hour)
		err := validatePresaleConfig(true, "CODE", nil, &general, &eventDate)
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("presale active with nil general sale returns ErrInvalidInput", func(t *testing.T) {
		start := time.Now().Add(1 * time.Hour)
		err := validatePresaleConfig(true, "CODE", &start, nil, &eventDate)
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("presale active with empty code returns ErrInvalidInput", func(t *testing.T) {
		start := time.Now().Add(1 * time.Hour)
		general := time.Now().Add(24 * time.Hour)
		err := validatePresaleConfig(true, "", &start, &general, &eventDate)
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("presale start after general sale returns ErrInvalidInput", func(t *testing.T) {
		start := time.Now().Add(24 * time.Hour)
		general := time.Now().Add(1 * time.Hour)
		err := validatePresaleConfig(true, "CODE", &start, &general, &eventDate)
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("general sale after event date returns ErrInvalidInput", func(t *testing.T) {
		start := time.Now().Add(1 * time.Hour)
		general := time.Now().Add(100 * time.Hour)
		err := validatePresaleConfig(true, "CODE", &start, &general, &eventDate)
		assert.ErrorIs(t, err, ErrInvalidInput)
	})

	t.Run("valid presale config returns nil", func(t *testing.T) {
		start := time.Now().Add(1 * time.Hour)
		general := time.Now().Add(24 * time.Hour)
		err := validatePresaleConfig(true, "CODE", &start, &general, &eventDate)
		assert.NoError(t, err)
	})
}
