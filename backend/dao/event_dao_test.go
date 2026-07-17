package db

import (
	"testing"
	"time"

	"backend/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestEvent(t *testing.T, dao *EventDAOImpl, title string, eventDate time.Time, status string, category string) *domain.Event {
	t.Helper()
	cat := category
	e := &domain.Event{
		Title:       title,
		EventDate:   eventDate,
		Capacity:    100,
		TicketsSold: 0,
		Price:       50.0,
		Status:      status,
		Category:    &cat,
		CreatedByID: 1,
	}
	require.NoError(t, dao.Create(e))
	return e
}

func TestEventDAO_Create(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewEventDAO(gdb)

	e := &domain.Event{
		Title:       "Concert",
		EventDate:   time.Now().Add(24 * time.Hour),
		Capacity:    200,
		TicketsSold: 0,
		Price:       100.0,
		Status:      "active",
		CreatedByID: 1,
	}
	err := dao.Create(e)
	require.NoError(t, err)
	assert.NotZero(t, e.ID)
}

func TestEventDAO_FindByID(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewEventDAO(gdb)

	t.Run("finds existing event", func(t *testing.T) {
		e := createTestEvent(t, dao, "FindMe", time.Now().Add(48*time.Hour), "active", "music")
		found, err := dao.FindByID(e.ID)
		require.NoError(t, err)
		assert.Equal(t, "FindMe", found.Title)
	})

	t.Run("returns error for non-existent event", func(t *testing.T) {
		_, err := dao.FindByID(99999)
		assert.Error(t, err)
	})
}

func TestEventDAO_FindAll(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewEventDAO(gdb)

	u := &domain.User{Name: "Admin", Email: "admin@test.com", PasswordHash: "h", Role: "admin"}
	require.NoError(t, gdb.Create(u).Error)

	futureDate := time.Now().Add(72 * time.Hour)
	createTestEvent(t, dao, "Event1", futureDate, "active", "music")
	createTestEvent(t, dao, "Event2", futureDate, "active", "sports")
	createTestEvent(t, dao, "Event3", futureDate, "active", "music")

	t.Run("returns all future events", func(t *testing.T) {
		events, err := dao.FindAll(domain.EventFilters{})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(events), 3)
	})

	t.Run("filters by category", func(t *testing.T) {
		events, err := dao.FindAll(domain.EventFilters{Category: "sports"})
		require.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, "Event2", events[0].Title)
	})

	t.Run("filters by date range", func(t *testing.T) {
		from := futureDate.Add(-1 * time.Hour)
		to := futureDate.Add(1 * time.Hour)
		events, err := dao.FindAll(domain.EventFilters{DateFrom: &from, DateTo: &to})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(events), 3)
	})

	t.Run("filters by price range", func(t *testing.T) {
		minP := 40.0
		maxP := 60.0
		events, err := dao.FindAll(domain.EventFilters{MinPrice: &minP, MaxPrice: &maxP})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(events), 3)
	})

	t.Run("hides past events", func(t *testing.T) {
		pastDate := time.Now().Add(-72 * time.Hour)
		createTestEvent(t, dao, "PastEvent", pastDate, "active", "music")
		events, err := dao.FindAll(domain.EventFilters{})
		require.NoError(t, err)
		for _, ev := range events {
			assert.False(t, ev.Title == "PastEvent", "past event should be hidden")
		}
	})
}

func TestEventDAO_Update(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewEventDAO(gdb)

	e := createTestEvent(t, dao, "Before", time.Now().Add(24*time.Hour), "active", "music")
	e.Title = "After"
	e.Price = 99.99
	err := dao.Update(e)
	require.NoError(t, err)

	found, err := dao.FindByID(e.ID)
	require.NoError(t, err)
	assert.Equal(t, "After", found.Title)
	assert.Equal(t, 99.99, found.Price)
}

func TestEventDAO_Delete(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewEventDAO(gdb)

	e := createTestEvent(t, dao, "ToDelete", time.Now().Add(24*time.Hour), "active", "music")
	err := dao.Delete(e.ID)
	require.NoError(t, err)

	_, err = dao.FindByID(e.ID)
	assert.Error(t, err)
}

func TestEventDAO_IncrementTicketsSold(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewEventDAO(gdb)

	e := createTestEvent(t, dao, "IncEvent", time.Now().Add(24*time.Hour), "active", "music")
	require.Equal(t, 0, e.TicketsSold)

	err := dao.IncrementTicketsSold(e.ID, 3)
	require.NoError(t, err)

	found, err := dao.FindByID(e.ID)
	require.NoError(t, err)
	assert.Equal(t, 3, found.TicketsSold)
}

func TestEventDAO_DecrementTicketsSold(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewEventDAO(gdb)

	e := createTestEvent(t, dao, "DecEvent", time.Now().Add(24*time.Hour), "active", "music")
	require.NoError(t, dao.IncrementTicketsSold(e.ID, 5))

	err := dao.DecrementTicketsSold(e.ID)
	require.NoError(t, err)

	found, err := dao.FindByID(e.ID)
	require.NoError(t, err)
	assert.Equal(t, 4, found.TicketsSold)
}

func TestEventDAO_DecrementTicketsSoldBelowZero(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewEventDAO(gdb)

	e := createTestEvent(t, dao, "DecZero", time.Now().Add(24*time.Hour), "active", "music")

	err := dao.DecrementTicketsSold(e.ID)
	require.NoError(t, err)

	found, err := dao.FindByID(e.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, found.TicketsSold)
}
