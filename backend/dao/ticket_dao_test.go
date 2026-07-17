package db

import (
	"testing"
	"time"

	"backend/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func seedTicketData(t *testing.T, gdb *gorm.DB) (*domain.User, *domain.Event) {
	t.Helper()
	u := &domain.User{Name: "TicketUser", Email: "ticketuser@test.com", PasswordHash: "h", Role: "client"}
	require.NoError(t, gdb.Create(u).Error)
	cat := "music"
	e := &domain.Event{
		Title:       "TicketEvent",
		EventDate:   time.Now().Add(72 * time.Hour),
		Capacity:    100,
		TicketsSold: 0,
		Price:       75.0,
		Status:      "active",
		Category:    &cat,
		CreatedByID: u.ID,
	}
	require.NoError(t, gdb.Create(e).Error)
	return u, e
}

func TestTicketDAO_Create(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewTicketDAO(gdb)
	u, e := seedTicketData(t, gdb)

	ticket := &domain.Ticket{
		UserID:        u.ID,
		EventID:       e.ID,
		Status:        "active",
		PurchasePrice: 75.0,
		PurchasedAt:   time.Now(),
	}
	err := dao.Create(ticket)
	require.NoError(t, err)
	assert.NotZero(t, ticket.ID)
}

func TestTicketDAO_FindByID(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewTicketDAO(gdb)
	u, e := seedTicketData(t, gdb)

	t.Run("finds existing ticket", func(t *testing.T) {
		ticket := &domain.Ticket{
			UserID: u.ID, EventID: e.ID, Status: "active",
			PurchasePrice: 75.0, PurchasedAt: time.Now(),
		}
		require.NoError(t, dao.Create(ticket))
		found, err := dao.FindByID(ticket.ID)
		require.NoError(t, err)
		assert.Equal(t, u.ID, found.UserID)
		assert.Equal(t, e.ID, found.EventID)
	})

	t.Run("returns error for non-existent ticket", func(t *testing.T) {
		_, err := dao.FindByID(99999)
		assert.Error(t, err)
	})
}

func TestTicketDAO_FindByUserID(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewTicketDAO(gdb)
	u, e := seedTicketData(t, gdb)

	t.Run("returns tickets for user", func(t *testing.T) {
		ticket := &domain.Ticket{
			UserID: u.ID, EventID: e.ID, Status: "active",
			PurchasePrice: 75.0, PurchasedAt: time.Now(),
		}
		require.NoError(t, dao.Create(ticket))
		tickets, err := dao.FindByUserID(u.ID)
		require.NoError(t, err)
		assert.Len(t, tickets, 1)
	})

	t.Run("returns empty list for user with no tickets", func(t *testing.T) {
		u2 := &domain.User{Name: "NoTickets", Email: "notickets@test.com", PasswordHash: "h", Role: "client"}
		require.NoError(t, gdb.Create(u2).Error)
		tickets, err := dao.FindByUserID(u2.ID)
		require.NoError(t, err)
		assert.Empty(t, tickets)
	})
}

func TestTicketDAO_FindActiveByEvent(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewTicketDAO(gdb)
	u, e := seedTicketData(t, gdb)

	t1 := &domain.Ticket{
		UserID: u.ID, EventID: e.ID, Status: "active",
		PurchasePrice: 75.0, PurchasedAt: time.Now(),
	}
	require.NoError(t, dao.Create(t1))

	tickets, err := dao.FindActiveByEvent(e.ID)
	require.NoError(t, err)
	assert.Len(t, tickets, 1)
}

func TestTicketDAO_CountActiveByEvent(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewTicketDAO(gdb)
	u, e := seedTicketData(t, gdb)

	for i := 0; i < 3; i++ {
		ticket := &domain.Ticket{
			UserID: u.ID, EventID: e.ID, Status: "active",
			PurchasePrice: 75.0, PurchasedAt: time.Now(),
		}
		require.NoError(t, dao.Create(ticket))
	}

	count, err := dao.CountActiveByEvent(e.ID)
	require.NoError(t, err)
	assert.Equal(t, 3, count)
}

func TestTicketDAO_CancelByEvent(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewTicketDAO(gdb)
	u, e := seedTicketData(t, gdb)

	t1 := &domain.Ticket{
		UserID: u.ID, EventID: e.ID, Status: "active",
		PurchasePrice: 75.0, PurchasedAt: time.Now(),
	}
	require.NoError(t, dao.Create(t1))

	err := dao.CancelByEvent(e.ID)
	require.NoError(t, err)

	count, err := dao.CountActiveByEvent(e.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	found, err := dao.FindByID(t1.ID)
	require.NoError(t, err)
	assert.Equal(t, "cancelled", found.Status)
}

func TestTicketDAO_Save(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewTicketDAO(gdb)
	u, e := seedTicketData(t, gdb)

	ticket := &domain.Ticket{
		UserID: u.ID, EventID: e.ID, Status: "active",
		PurchasePrice: 75.0, PurchasedAt: time.Now(),
	}
	require.NoError(t, dao.Create(ticket))

	ticket.Status = "transferred"
	err := dao.Save(ticket)
	require.NoError(t, err)

	found, err := dao.FindByID(ticket.ID)
	require.NoError(t, err)
	assert.Equal(t, "transferred", found.Status)
}

func TestTicketDAO_WithTransaction(t *testing.T) {
	gdb := setupTestDB(t)
	dao := NewTicketDAO(gdb)
	u, e := seedTicketData(t, gdb)

	t.Run("successful transaction", func(t *testing.T) {
		err := dao.WithTransaction(func(tx TxContext) error {
			ticket := &domain.Ticket{
				UserID: u.ID, EventID: e.ID, Status: "active",
				PurchasePrice: 75.0, PurchasedAt: time.Now(),
			}
			return dao.Create(ticket)
		})
		require.NoError(t, err)

		tickets, err := dao.FindByUserID(u.ID)
		require.NoError(t, err)
		assert.Len(t, tickets, 1)
	})

	t.Run("transaction returns error", func(t *testing.T) {
		err := dao.WithTransaction(func(tx TxContext) error {
			return assert.AnError
		})
		assert.ErrorIs(t, err, assert.AnError)
	})
}
