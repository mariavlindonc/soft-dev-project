package controllers

import (
	"net/http"

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

func (h *AdminController) GetReports(c *gin.Context) {
	report, err := h.reportService.GetGlobalReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toGlobalReportResponse(report))
}
