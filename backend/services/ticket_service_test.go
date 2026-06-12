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

func TestPurchase(t *testing.T) {
	activeEvent := &domain.Event{ID: 1, Status: "active", Capacity: 100, Price: 50}

	t.Run("success purchases ticket", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		ticketDAO := new(MockTicketDAO)
		userDAO := new(MockUserDAO)
		email := new(MockEmailClient)
		svc := NewTicketService(ticketDAO, eventDAO, userDAO, email)

		eventDAO.On("FindByID", uint(1)).Return(activeEvent, nil)
		ticketDAO.On("CountActiveByEvent", uint(1)).Return(0, nil)
		ticketDAO.On("WithTransaction", mock.Anything).Return(nil)
		// The transaction calls Create and IncrementTicketsSold
		ticketDAO.On("Create", mock.AnythingOfType("*domain.Ticket")).Return(nil)
		eventDAO.On("IncrementTicketsSold", uint(1), 1).Return(nil)
		userDAO.On("FindByID", uint(10)).Return(&domain.User{ID: 10, Email: "buyer@test.com"}, nil)
		email.On("SendPurchaseConfirmation", "buyer@test.com", mock.Anything).Return(nil)

		tickets, err := svc.Purchase(10, PurchaseInput{EventID: 1, Quantity: 1})
		require.NoError(t, err)
		require.Len(t, tickets, 1)
		assert.Equal(t, uint(1), tickets[0].EventID)
		assert.Equal(t, "active", tickets[0].Status)
		assert.Equal(t, 50.0, tickets[0].PurchasePrice)
	})

	t.Run("cancelled event returns ErrEventCancelled", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		svc := NewTicketService(new(MockTicketDAO), eventDAO, new(MockUserDAO), new(MockEmailClient))

		eventDAO.On("FindByID", uint(1)).Return(&domain.Event{ID: 1, Status: "cancelled"}, nil)

		_, err := svc.Purchase(10, PurchaseInput{EventID: 1})
		assert.ErrorIs(t, err, ErrEventCancelled)
	})

	t.Run("full capacity returns ErrNoCapacity", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		ticketDAO := new(MockTicketDAO)
		svc := NewTicketService(ticketDAO, eventDAO, new(MockUserDAO), new(MockEmailClient))

		fullEvent := &domain.Event{ID: 2, Status: "active", Capacity: 100}
		eventDAO.On("FindByID", uint(2)).Return(fullEvent, nil)
		ticketDAO.On("CountActiveByEvent", uint(2)).Return(100, nil)

		_, err := svc.Purchase(10, PurchaseInput{EventID: 2})
		assert.ErrorIs(t, err, ErrNoCapacity)
	})
}

func TestPurchasePresale(t *testing.T) {
	now := time.Now()
	presaleStart := now.Add(-2 * time.Hour)
	generalSale := now.Add(2 * time.Hour)
	code := "PRE123"

	presaleEvent := &domain.Event{
		ID:               1,
		Status:           "presale",
		Capacity:         100,
		PresaleActive:    true,
		PresaleCode:      &code,
		PresaleStartDate: &presaleStart,
		GeneralSaleDate:  &generalSale,
		Price:            50,
	}

	t.Run("presale with correct code succeeds", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		ticketDAO := new(MockTicketDAO)
		userDAO := new(MockUserDAO)
		email := new(MockEmailClient)
		svc := NewTicketService(ticketDAO, eventDAO, userDAO, email)

		eventDAO.On("FindByID", uint(1)).Return(presaleEvent, nil)
		ticketDAO.On("CountActiveByEvent", uint(1)).Return(0, nil)
		ticketDAO.On("WithTransaction", mock.Anything).Return(nil)
		ticketDAO.On("Create", mock.AnythingOfType("*domain.Ticket")).Return(nil)
		eventDAO.On("IncrementTicketsSold", uint(1), 1).Return(nil)
		userDAO.On("FindByID", uint(10)).Return(&domain.User{ID: 10, Email: "u@test.com"}, nil)
		email.On("SendPurchaseConfirmation", "u@test.com", mock.Anything).Return(nil)

		_, err := svc.Purchase(10, PurchaseInput{EventID: 1, Quantity: 1, PresaleCode: "PRE123"})
		require.NoError(t, err)
	})

	t.Run("presale without code returns ErrPresaleCodeRequired", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		ticketDAO := new(MockTicketDAO)
		svc := NewTicketService(ticketDAO, eventDAO, new(MockUserDAO), new(MockEmailClient))

		eventDAO.On("FindByID", uint(1)).Return(presaleEvent, nil)

		_, err := svc.Purchase(10, PurchaseInput{EventID: 1})
		assert.ErrorIs(t, err, ErrPresaleCodeRequired)
	})

	t.Run("presale with wrong code returns ErrInvalidPresaleCode", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		ticketDAO := new(MockTicketDAO)
		svc := NewTicketService(ticketDAO, eventDAO, new(MockUserDAO), new(MockEmailClient))

		eventDAO.On("FindByID", uint(1)).Return(presaleEvent, nil)

		_, err := svc.Purchase(10, PurchaseInput{EventID: 1, PresaleCode: "WRONG"})
		assert.ErrorIs(t, err, ErrInvalidPresaleCode)
	})
}

func TestCancelTicket(t *testing.T) {
	userID := uint(10)

	t.Run("cancel own ticket succeeds", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		ticketDAO := new(MockTicketDAO)
		userDAO := new(MockUserDAO)
		email := new(MockEmailClient)
		svc := NewTicketService(ticketDAO, eventDAO, userDAO, email)

		ticket := &domain.Ticket{ID: 1, UserID: userID, Status: "active", EventID: 5}
		ticketDAO.On("FindByID", uint(1)).Return(ticket, nil)
		ticketDAO.On("WithTransaction", mock.Anything).Return(nil)
		ticketDAO.On("Save", mock.MatchedBy(func(t *domain.Ticket) bool {
			return t.Status == "cancelled" && t.CancelledAt != nil
		})).Return(nil)
		eventDAO.On("DecrementTicketsSold", uint(5)).Return(nil)
		userDAO.On("FindByID", userID).Return(&domain.User{ID: userID, Email: "u@test.com"}, nil)
		eventDAO.On("FindByID", uint(5)).Return(&domain.Event{ID: 5, Title: "Concert"}, nil)
		email.On("SendCancellationNotice", "u@test.com", mock.Anything).Return(nil)

		err := svc.Cancel(1, userID)
		require.NoError(t, err)
	})

	t.Run("cancel someone else's ticket returns ErrNotFound", func(t *testing.T) {
		ticketDAO := new(MockTicketDAO)
		svc := NewTicketService(ticketDAO, new(MockEventDAO), new(MockUserDAO), new(MockEmailClient))

		ticketDAO.On("FindByID", uint(1)).Return(&domain.Ticket{ID: 1, UserID: 99}, nil)

		err := svc.Cancel(1, userID)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("cancel already cancelled ticket returns ErrAlreadyCancelled", func(t *testing.T) {
		ticketDAO := new(MockTicketDAO)
		svc := NewTicketService(ticketDAO, new(MockEventDAO), new(MockUserDAO), new(MockEmailClient))

		ticketDAO.On("FindByID", uint(1)).Return(&domain.Ticket{ID: 1, UserID: userID, Status: "cancelled"}, nil)

		err := svc.Cancel(1, userID)
		assert.ErrorIs(t, err, ErrAlreadyCancelled)
	})
}

func TestTransfer(t *testing.T) {
	fromUserID := uint(10)
	targetUser := &domain.User{ID: 20, Email: "target@test.com", Name: "Target"}

	t.Run("transfer to another user succeeds", func(t *testing.T) {
		eventDAO := new(MockEventDAO)
		ticketDAO := new(MockTicketDAO)
		userDAO := new(MockUserDAO)
		email := new(MockEmailClient)
		svc := NewTicketService(ticketDAO, eventDAO, userDAO, email)

		ticket := &domain.Ticket{ID: 1, UserID: fromUserID, Status: "active", EventID: 5}
		ticketDAO.On("FindByID", uint(1)).Return(ticket, nil)
		userDAO.On("FindByEmail", "target@test.com").Return(targetUser, nil)
		ticketDAO.On("WithTransaction", mock.Anything).Return(nil)
		ticketDAO.On("Save", mock.MatchedBy(func(t *domain.Ticket) bool {
			return t.Status == "transferred" && *t.TransferredToID == 20
		})).Return(nil)
		userDAO.On("FindByID", fromUserID).Return(&domain.User{ID: fromUserID, Email: "from@test.com"}, nil)
		eventDAO.On("FindByID", uint(5)).Return(&domain.Event{ID: 5, Title: "Concert"}, nil)
		email.On("SendTransferNotice", "from@test.com", "target@test.com", mock.Anything).Return(nil)

		err := svc.Transfer(1, fromUserID, TransferInput{ToUserEmail: "target@test.com"})
		require.NoError(t, err)
	})

	t.Run("transfer to self returns ErrInvalidTransfer", func(t *testing.T) {
		userDAO := new(MockUserDAO)
		ticketDAO := new(MockTicketDAO)
		svc := NewTicketService(ticketDAO, new(MockEventDAO), userDAO, new(MockEmailClient))

		ticketDAO.On("FindByID", uint(1)).Return(&domain.Ticket{ID: 1, UserID: fromUserID, Status: "active"}, nil)
		userDAO.On("FindByEmail", "self@test.com").Return(&domain.User{ID: fromUserID, Email: "self@test.com"}, nil)

		err := svc.Transfer(1, fromUserID, TransferInput{ToUserEmail: "self@test.com"})
		assert.ErrorIs(t, err, ErrInvalidTransfer)
	})
}

func TestGetByUser(t *testing.T) {
	ticketDAO := new(MockTicketDAO)
	svc := NewTicketService(ticketDAO, new(MockEventDAO), new(MockUserDAO), new(MockEmailClient))

	expected := []domain.Ticket{{ID: 1, UserID: 10}, {ID: 2, UserID: 10}}
	ticketDAO.On("FindByUserID", uint(10)).Return(expected, nil)

	tickets, err := svc.GetByUser(10)
	require.NoError(t, err)
	assert.Len(t, tickets, 2)
}

func TestPurchaseNotFound(t *testing.T) {
	eventDAO := new(MockEventDAO)
	svc := NewTicketService(new(MockTicketDAO), eventDAO, new(MockUserDAO), new(MockEmailClient))

	eventDAO.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	_, err := svc.Purchase(10, PurchaseInput{EventID: 999})
	assert.ErrorIs(t, err, ErrNotFound)
}
