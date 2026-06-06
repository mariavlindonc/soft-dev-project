package controllers

import (
	"net/http"
	"strconv"

	"backend/services"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	reportService services.ReportServiceInterface
}

func NewAdminController(s services.ReportServiceInterface) *AdminController {
	return &AdminController{reportService: s}
}

type eventSummary struct {
	EventID     uint    `json:"event_id"`
	Title       string  `json:"title"`
	Capacity    int     `json:"capacity"`
	TicketsSold int     `json:"tickets_sold"`
	Occupancy   float64 `json:"occupancy"`
}

type globalReportResponse struct {
	TotalEvents      int            `json:"total_events"`
	TotalTicketsSold int            `json:"total_tickets_sold"`
	Events           []eventSummary `json:"events"`
}

func toGlobalReportResponse(report *services.GlobalReport) globalReportResponse {
	resp := globalReportResponse{
		TotalEvents:      report.TotalEvents,
		TotalTicketsSold: report.TotalTicketsSold,
		Events:           make([]eventSummary, 0, len(report.EventReports)),
	}

	for _, e := range report.EventReports {
		resp.Events = append(resp.Events, eventSummary{
			EventID:     e.EventID,
			Title:       e.EventTitle,
			Capacity:    e.TotalCapacity,
			TicketsSold: e.TicketsSold,
			Occupancy:   e.Occupancy,
		})
	}

	return resp
}

type buyerInfoResponse struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type eventReportResponse struct {
	EventID       uint               `json:"event_id"`
	Title         string             `json:"title"`
	Capacity      int                `json:"capacity"`
	TicketsSold   int                `json:"tickets_sold"`
	Occupancy     float64            `json:"occupancy"`
	Buyers        []buyerInfoResponse `json:"buyers"`
}

func toEventReportResponse(report *services.EventReport) eventReportResponse {
	resp := eventReportResponse{
		EventID:     report.EventID,
		Title:       report.EventTitle,
		Capacity:    report.TotalCapacity,
		TicketsSold: report.TicketsSold,
		Occupancy:   report.Occupancy,
		Buyers:      make([]buyerInfoResponse, 0, len(report.Buyers)),
	}
	for _, b := range report.Buyers {
		resp.Buyers = append(resp.Buyers, buyerInfoResponse{
			UserID: b.UserID,
			Name:   b.Name,
			Email:  b.Email,
		})
	}
	return resp
}

func (h *AdminController) GetReports(c *gin.Context) {
	report, err := h.reportService.GetGlobalReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toGlobalReportResponse(report))
}

func (h *AdminController) GetEventReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	report, err := h.reportService.GetEventReport(uint(id))
	if err != nil {
		if isNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toEventReportResponse(report))
}
