package services

import (
	"fmt"

	db "backend/dao"
	"backend/domain"
)

// ---------------------------------------------------------------------------
// Interface
// ---------------------------------------------------------------------------

type ReportServiceInterface interface {
	GetEventReport(eventID uint) (*EventReport, error)
	GetGlobalReport() (*GlobalReport, error)
}

// ---------------------------------------------------------------------------
// Output types
// ---------------------------------------------------------------------------

type EventReport struct {
	EventID       uint
	EventTitle    string
	TotalCapacity int
	TicketsSold   int
	Occupancy     float64
	Buyers        []BuyerInfo
}

type BuyerInfo struct {
	UserID uint
	Name   string
	Email  string
}

type GlobalReport struct {
	TotalEvents      int
	TotalTicketsSold int
	EventReports     []EventReport
}

// ---------------------------------------------------------------------------
// Concrete service
// ---------------------------------------------------------------------------

type ReportService struct {
	eventDAO  db.EventDAO
	ticketDAO db.TicketDAO
	userDAO   db.UserDAO
}

func NewReportService(
	eventDAO db.EventDAO,
	ticketDAO db.TicketDAO,
	userDAO db.UserDAO,
) *ReportService {
	return &ReportService{
		eventDAO:  eventDAO,
		ticketDAO: ticketDAO,
		userDAO:   userDAO,
	}
}

func (s *ReportService) GetEventReport(eventID uint) (*EventReport, error) {
	event, err := s.eventDAO.FindByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("event %d: %w", eventID, ErrNotFound)
	}

	ticketsSold, err := s.ticketDAO.CountActiveByEvent(eventID)
	if err != nil {
		return nil, fmt.Errorf("count tickets: %w", err)
	}

	occupancy := 0.0
	if event.Capacity > 0 {
		occupancy = float64(ticketsSold) / float64(event.Capacity) * 100
	}

	report := &EventReport{
		EventID:       event.ID,
		EventTitle:    event.Title,
		TotalCapacity: event.Capacity,
		TicketsSold:   ticketsSold,
		Occupancy:     occupancy,
	}

	activeTickets, err := s.ticketDAO.FindActiveByEvent(eventID)
	if err == nil {
		for _, t := range activeTickets {
			user, uerr := s.userDAO.FindByID(t.UserID)
			if uerr == nil {
				report.Buyers = append(report.Buyers, BuyerInfo{
					UserID: user.ID,
					Name:   user.Name,
					Email:  user.Email,
				})
			}
		}
	}

	return report, nil
}

func (s *ReportService) GetGlobalReport() (*GlobalReport, error) {
	// Fetch all events (empty filters = no filter).
	events, err := s.eventDAO.FindAll(domain.EventFilters{})
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}

	report := &GlobalReport{
		TotalEvents: len(events),
	}

	for _, event := range events {
		ticketsSold, err := s.ticketDAO.CountActiveByEvent(event.ID)
		if err != nil {
			return nil, fmt.Errorf("count tickets for event %d: %w", event.ID, err)
		}

		report.TotalTicketsSold += ticketsSold

		occupancy := 0.0
		if event.Capacity > 0 {
			occupancy = float64(ticketsSold) / float64(event.Capacity) * 100
		}

		report.EventReports = append(report.EventReports, EventReport{
			EventID:       event.ID,
			EventTitle:    event.Title,
			TotalCapacity: event.Capacity,
			TicketsSold:   ticketsSold,
			Occupancy:     occupancy,
		})
	}

	return report, nil
}

var _ ReportServiceInterface = (*ReportService)(nil)
