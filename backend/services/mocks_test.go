package services

import (
	"backend/clients"
	db "backend/dao"
	"backend/domain"

	"github.com/stretchr/testify/mock"
)

type MockUserDAO struct {
	mock.Mock
}

func (m *MockUserDAO) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserDAO) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserDAO) FindByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

var _ db.UserDAO = (*MockUserDAO)(nil)

type MockEventDAO struct {
	mock.Mock
}

func (m *MockEventDAO) FindAll(filters domain.EventFilters) ([]domain.Event, error) {
	args := m.Called(filters)
	return args.Get(0).([]domain.Event), args.Error(1)
}

func (m *MockEventDAO) FindByID(id uint) (*domain.Event, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Event), args.Error(1)
}

func (m *MockEventDAO) Create(event *domain.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventDAO) Update(event *domain.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventDAO) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEventDAO) IncrementTicketsSold(eventID uint) error {
	args := m.Called(eventID)
	return args.Error(0)
}

func (m *MockEventDAO) DecrementTicketsSold(eventID uint) error {
	args := m.Called(eventID)
	return args.Error(0)
}

var _ db.EventDAO = (*MockEventDAO)(nil)

type MockTicketDAO struct {
	mock.Mock
}

func (m *MockTicketDAO) Create(ticket *domain.Ticket) error {
	args := m.Called(ticket)
	return args.Error(0)
}

func (m *MockTicketDAO) FindByID(id uint) (*domain.Ticket, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ticket), args.Error(1)
}

func (m *MockTicketDAO) FindByUserID(userID uint) ([]domain.Ticket, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.Ticket), args.Error(1)
}

func (m *MockTicketDAO) FindActiveByEvent(eventID uint) ([]domain.Ticket, error) {
	args := m.Called(eventID)
	return args.Get(0).([]domain.Ticket), args.Error(1)
}

func (m *MockTicketDAO) CountActiveByEvent(eventID uint) (int, error) {
	args := m.Called(eventID)
	return args.Int(0), args.Error(1)
}

func (m *MockTicketDAO) CancelByEvent(eventID uint) error {
	args := m.Called(eventID)
	return args.Error(0)
}

func (m *MockTicketDAO) Save(ticket *domain.Ticket) error {
	args := m.Called(ticket)
	return args.Error(0)
}

func (m *MockTicketDAO) WithTransaction(fn func(db.TxContext) error) error {
	args := m.Called(fn)
	fn(nil)
	return args.Error(0)
}

var _ db.TicketDAO = (*MockTicketDAO)(nil)

type MockEmailClient struct {
	mock.Mock
}

func (m *MockEmailClient) SendPurchaseConfirmation(to string, ticket clients.TicketInfo) error {
	args := m.Called(to, ticket)
	return args.Error(0)
}

func (m *MockEmailClient) SendCancellationNotice(to string, ticket clients.TicketInfo) error {
	args := m.Called(to, ticket)
	return args.Error(0)
}

func (m *MockEmailClient) SendTransferNotice(from, to string, ticket clients.TicketInfo) error {
	args := m.Called(from, to, ticket)
	return args.Error(0)
}

var _ clients.EmailClient = (*MockEmailClient)(nil)
