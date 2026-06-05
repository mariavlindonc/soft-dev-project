package controllers

import (
	"net/http"

	"backend/services"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	adminService services.AdminService
}

func NewAdminController(s services.AdminService) *AdminController {
	return &AdminController{adminService: s}
}

type eventSummary struct {
	EventID     uint    `json:"event_id"`
	Title       string  `json:"title"`
	Capacity    int     `json:"capacity"`
	TicketsSold int     `json:"tickets_sold"`
	Occupancy   float64 `json:"occupancy"`
}

type reportResponse struct {
	TotalCapacity int            `json:"total_capacity"`
	TicketsSold   int            `json:"tickets_sold"`
	OccupancyRate float64        `json:"occupancy_rate"`
	Events        []eventSummary `json:"events"`
}

func toReportResponse(report *services.Report) reportResponse {
	resp := reportResponse{
		TotalCapacity: report.TotalCapacity,
		TicketsSold:   report.TicketsSold,
		Events:        make([]eventSummary, 0, len(report.Events)),
	}

	if resp.TotalCapacity > 0 {
		resp.OccupancyRate = float64(resp.TicketsSold) / float64(resp.TotalCapacity) * 100
	}

	for _, e := range report.Events {
		occ := 0.0
		if e.Capacity > 0 {
			occ = float64(e.TicketsSold) / float64(e.Capacity) * 100
		}
		resp.Events = append(resp.Events, eventSummary{
			EventID:     e.EventID,
			Title:       e.Title,
			Capacity:    e.Capacity,
			TicketsSold: e.TicketsSold,
			Occupancy:   occ,
		})
	}

	return resp
}

func (h *AdminController) GetReports(c *gin.Context) {
	report, err := h.adminService.GetReports()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toReportResponse(report))
}
