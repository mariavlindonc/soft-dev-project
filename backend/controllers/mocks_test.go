package controllers

import (
	"backend/domain"
	"backend/services"

	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(input services.RegisterInput) (*domain.User, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockAuthService) Login(input services.LoginInput) (string, error) {
	args := m.Called(input)
	return args.String(0), args.Error(1)
}

var _ services.AuthServiceInterface = (*MockAuthService)(nil)

type MockEventService struct {
	mock.Mock
}

func (m *MockEventService) GetAll(filters domain.EventFilters) ([]domain.Event, error) {
	args := m.Called(filters)
	return args.Get(0).([]domain.Event), args.Error(1)
}

func (m *MockEventService) GetByID(id uint) (*domain.Event, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Event), args.Error(1)
}

func (m *MockEventService) Create(input services.CreateEventInput) (*domain.Event, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Event), args.Error(1)
}

func (m *MockEventService) Update(id uint, input services.UpdateEventInput) (*domain.Event, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Event), args.Error(1)
}

func (m *MockEventService) Cancel(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

var _ services.EventServiceInterface = (*MockEventService)(nil)

type MockTicketService struct {
	mock.Mock
}

func (m *MockTicketService) Purchase(userID uint, input services.PurchaseInput) (*domain.Ticket, error) {
	args := m.Called(userID, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Ticket), args.Error(1)
}

func (m *MockTicketService) GetByUser(userID uint) ([]domain.Ticket, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.Ticket), args.Error(1)
}

func (m *MockTicketService) Cancel(ticketID uint, requestingUserID uint) error {
	args := m.Called(ticketID, requestingUserID)
	return args.Error(0)
}

func (m *MockTicketService) Transfer(ticketID uint, fromUserID uint, input services.TransferInput) error {
	args := m.Called(ticketID, fromUserID, input)
	return args.Error(0)
}

var _ services.TicketServiceInterface = (*MockTicketService)(nil)

type MockReportService struct {
	mock.Mock
}

func (m *MockReportService) GetEventReport(eventID uint) (*services.EventReport, error) {
	args := m.Called(eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.EventReport), args.Error(1)
}

func (m *MockReportService) GetGlobalReport() (*services.GlobalReport, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.GlobalReport), args.Error(1)
}

var _ services.ReportServiceInterface = (*MockReportService)(nil)
