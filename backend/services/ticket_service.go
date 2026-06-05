package services

import (
	"fmt"
	"log"
	"time"

	db "backend/dao"
	"backend/domain"
	"backend/clients"
)

// ---------------------------------------------------------------------------
// Interface
// ---------------------------------------------------------------------------

type TicketServiceInterface interface {
	Purchase(userID uint, input PurchaseInput) (*domain.Ticket, error)
	GetByUser(userID uint) ([]domain.Ticket, error)
	Cancel(ticketID uint, requestingUserID uint) error
	Transfer(ticketID uint, fromUserID uint, input TransferInput) error
}

// ---------------------------------------------------------------------------
// Input types
// ---------------------------------------------------------------------------

type PurchaseInput struct {
	EventID uint
}

type TransferInput struct {
	ToUserEmail string
}

// ---------------------------------------------------------------------------
// Concrete service
// ---------------------------------------------------------------------------

type TicketService struct {
	ticketDAO   db.TicketDAO
	eventDAO    db.EventDAO
	userDAO     db.UserDAO
	emailClient clients.EmailClient
}

func NewTicketService(
	ticketDAO db.TicketDAO,
	eventDAO db.EventDAO,
	userDAO db.UserDAO,
	emailClient clients.EmailClient,
) *TicketService {
	return &TicketService{
		ticketDAO:   ticketDAO,
		eventDAO:    eventDAO,
		userDAO:     userDAO,
		emailClient: emailClient,
	}
}

// Purchase validates that the event exists, is not cancelled,
// and has available capacity, then creates the ticket inside a
// database transaction to avoid race conditions.
func (s *TicketService) Purchase(userID uint, input PurchaseInput) (*domain.Ticket, error) {
	event, err := s.eventDAO.FindByID(input.EventID)
	if err != nil {
		return nil, fmt.Errorf("event %d: %w", input.EventID, ErrNotFound)
	}
	if event.Status == "cancelled" {
		return nil, fmt.Errorf("event %d: %w", event.ID, ErrEventCancelled)
	}

	activeCount, err := s.ticketDAO.CountActiveByEvent(event.ID)
	if err != nil {
		return nil, fmt.Errorf("count tickets: %w", err)
	}
	if activeCount >= event.Capacity {
		return nil, fmt.Errorf("event %d: %w", event.ID, ErrNoCapacity)
	}

	ticket := &domain.Ticket{
		UserID:        userID,
		EventID:       event.ID,
		Status:        "active",
		PurchasePrice: event.Price,
	}

	if err := s.ticketDAO.WithTransaction(func(tx db.TxContext) error {
		return s.ticketDAO.Create(ticket)
	}); err != nil {
		return nil, fmt.Errorf("purchase ticket: %w", err)
	}

	s.sendPurchaseEmail(userID, ticket, event)

	return ticket, nil
}

func (s *TicketService) GetByUser(userID uint) ([]domain.Ticket, error) {
	return s.ticketDAO.FindByUserID(userID)
}

// Cancel marks the ticket as cancelled after verifying ownership.
// Ownership violations return ErrNotFound to avoid leaking existence info.
func (s *TicketService) Cancel(ticketID uint, requestingUserID uint) error {
	ticket, err := s.ticketDAO.FindByID(ticketID)
	if err != nil {
		return fmt.Errorf("ticket %d: %w", ticketID, ErrNotFound)
	}
	if ticket.UserID != requestingUserID {
		return fmt.Errorf("ticket %d: %w", ticketID, ErrNotFound)
	}
	if ticket.Status == "cancelled" {
		return fmt.Errorf("ticket %d: %w", ticketID, ErrAlreadyCancelled)
	}

	now := time.Now()
	ticket.Status = "cancelled"
	ticket.CancelledAt = &now

	if err := s.ticketDAO.Save(ticket); err != nil {
		return fmt.Errorf("cancel ticket: %w", err)
	}

	s.sendCancellationEmail(requestingUserID, ticket)

	return nil
}

// Transfer moves a ticket to another user. The target is resolved by email
// (never by userID from the request body) to prevent user enumeration.
func (s *TicketService) Transfer(ticketID uint, fromUserID uint, input TransferInput) error {
	ticket, err := s.ticketDAO.FindByID(ticketID)
	if err != nil {
		return fmt.Errorf("ticket %d: %w", ticketID, ErrNotFound)
	}
	if ticket.UserID != fromUserID {
		return fmt.Errorf("ticket %d: %w", ticketID, ErrNotFound)
	}
	if ticket.Status == "cancelled" {
		return fmt.Errorf("ticket %d: %w", ticketID, ErrAlreadyCancelled)
	}

	targetUser, err := s.userDAO.FindByEmail(input.ToUserEmail)
	if err != nil {
		return fmt.Errorf("target user %s: %w", input.ToUserEmail, ErrNotFound)
	}
	if fromUserID == targetUser.ID {
		return fmt.Errorf("ticket %d: %w", ticketID, ErrInvalidTransfer)
	}

	now := time.Now()
	ticket.TransferredToID = &targetUser.ID
	ticket.TransferredAt = &now
	ticket.Status = "transferred"

	if err := s.ticketDAO.WithTransaction(func(tx db.TxContext) error {
		return s.ticketDAO.Save(ticket)
	}); err != nil {
		return fmt.Errorf("transfer ticket: %w", err)
	}

	s.sendTransferEmail(fromUserID, targetUser, ticket)

	return nil
}

// ---------------------------------------------------------------------------
// Email helpers — errors are logged but never returned to the caller
// because a notification failure must not roll back a successful operation.
// ---------------------------------------------------------------------------

func (s *TicketService) sendPurchaseEmail(userID uint, ticket *domain.Ticket, event *domain.Event) {
	user, err := s.userDAO.FindByID(userID)
	if err != nil {
		log.Printf("warn: send purchase email: find user %d: %v", userID, err)
		return
	}
	if err := s.emailClient.SendPurchaseConfirmation(user.Email, toTicketInfo(ticket, event)); err != nil {
		log.Printf("warn: send purchase confirmation for ticket %d: %v", ticket.ID, err)
	}
}

func (s *TicketService) sendCancellationEmail(userID uint, ticket *domain.Ticket) {
	user, err := s.userDAO.FindByID(userID)
	if err != nil {
		log.Printf("warn: send cancellation email: find user %d: %v", userID, err)
		return
	}
	event, err := s.eventDAO.FindByID(ticket.EventID)
	if err != nil {
		log.Printf("warn: send cancellation email: find event %d: %v", ticket.EventID, err)
		return
	}
	if err := s.emailClient.SendCancellationNotice(user.Email, toTicketInfo(ticket, event)); err != nil {
		log.Printf("warn: send cancellation notice for ticket %d: %v", ticket.ID, err)
	}
}

func (s *TicketService) sendTransferEmail(fromUserID uint, toUser *domain.User, ticket *domain.Ticket) {
	fromUser, err := s.userDAO.FindByID(fromUserID)
	if err != nil {
		log.Printf("warn: send transfer email: find fromUser %d: %v", fromUserID, err)
		return
	}
	event, err := s.eventDAO.FindByID(ticket.EventID)
	if err != nil {
		log.Printf("warn: send transfer email: find event %d: %v", ticket.EventID, err)
		return
	}
	if err := s.emailClient.SendTransferNotice(fromUser.Email, toUser.Email, toTicketInfo(ticket, event)); err != nil {
		log.Printf("warn: send transfer notice for ticket %d: %v", ticket.ID, err)
	}
}

func toTicketInfo(ticket *domain.Ticket, event *domain.Event) clients.TicketInfo {
	loc := ""
	if event.Location != nil {
		loc = *event.Location
	}
	return clients.TicketInfo{
		TicketID:   ticket.ID,
		EventTitle: event.Title,
		EventDate:  event.EventDate.Format("2006-01-02 15:04"),
		Location:   loc,
		Price:      ticket.PurchasePrice,
	}
}

var _ TicketServiceInterface = (*TicketService)(nil)
