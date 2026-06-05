package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"backend/domain"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	eventService services.EventService
}

func NewEventController(s services.EventService) *EventController {
	return &EventController{eventService: s}
}

type createEventRequest struct {
	Title            string  `json:"title"              binding:"required"`
	Description      string  `json:"description"`
	ImageURL         string  `json:"image_url"`
	Category         string  `json:"category"`
	Location         string  `json:"location"`
	EventDate        string  `json:"event_date"         binding:"required"`
	DurationMinutes  int     `json:"duration_minutes"`
	Capacity         int     `json:"capacity"           binding:"required,min=1"`
	Price            float64 `json:"price"              binding:"required,min=0"`
	PresaleCode      string  `json:"presale_code"`
	PresaleStartDate string  `json:"presale_start_date"`
	GeneralSaleDate  string  `json:"general_sale_date"`
}

type updateEventRequest struct {
	Title            string   `json:"title"`
	Description      *string  `json:"description"`
	ImageURL         *string  `json:"image_url"`
	Category         *string  `json:"category"`
	Location         *string  `json:"location"`
	EventDate        string   `json:"event_date"`
	DurationMinutes  *int     `json:"duration_minutes"`
	Capacity         *int     `json:"capacity"           binding:"omitempty,min=1"`
	Price            *float64 `json:"price"              binding:"omitempty,min=0"`
	PresaleCode      *string  `json:"presale_code"`
	PresaleStartDate *string  `json:"presale_start_date"`
	GeneralSaleDate  *string  `json:"general_sale_date"`
	Status           *string  `json:"status"             binding:"omitempty,oneof=active presale sold_out cancelled"`
}

type eventResponse struct {
	ID              uint    `json:"id"`
	Title           string  `json:"title"`
	Description     *string `json:"description,omitempty"`
	ImageURL        *string `json:"image_url,omitempty"`
	Category        *string `json:"category,omitempty"`
	Location        *string `json:"location,omitempty"`
	EventDate       string  `json:"event_date"`
	DurationMinutes int     `json:"duration_minutes"`
	Capacity        int     `json:"capacity"`
	TicketsSold     int     `json:"tickets_sold"`
	Price           float64 `json:"price"`
	Status          string  `json:"status"`
	CreatedByID     uint    `json:"created_by_id"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

func toEventResponse(e *domain.Event) eventResponse {
	return eventResponse{
		ID:              e.ID,
		Title:           e.Title,
		Description:     e.Description,
		ImageURL:        e.ImageURL,
		Category:        e.Category,
		Location:        e.Location,
		EventDate:       e.EventDate.Format(time.RFC3339),
		DurationMinutes: e.DurationMinutes,
		Capacity:        e.Capacity,
		TicketsSold:     e.TicketsSold,
		Price:           e.Price,
		Status:          e.Status,
		CreatedByID:     e.CreatedByID,
		CreatedAt:       e.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       e.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *EventController) GetAll(c *gin.Context) {
	events, err := h.eventService.GetAll(c.Query("category"), c.Query("date"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]eventResponse, 0, len(events))
	for _, e := range events {
		resp = append(resp, toEventResponse(&e))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *EventController) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	event, err := h.eventService.GetByID(uint(id))
	if err != nil {
		if isNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toEventResponse(event))
}

func (h *EventController) Create(c *gin.Context) {
	var req createEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := toCreateEvent(&req, c.GetUint("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.eventService.Create(event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toEventResponse(created))
}

func toCreateEvent(req *createEventRequest, userID uint) (*domain.Event, error) {
	eventDate, err := time.Parse(time.RFC3339, req.EventDate)
	if err != nil {
		return nil, fmt.Errorf("invalid event_date, use RFC3339")
	}

	event := &domain.Event{
		Title:           req.Title,
		Description:     optionalString(req.Description),
		ImageURL:        optionalString(req.ImageURL),
		Category:        optionalString(req.Category),
		Location:        optionalString(req.Location),
		EventDate:       eventDate,
		DurationMinutes: req.DurationMinutes,
		Capacity:        req.Capacity,
		Price:           req.Price,
		CreatedByID:     userID,
	}

	if req.PresaleCode != "" {
		event.PresaleCode = &req.PresaleCode
	}
	if req.PresaleStartDate != "" {
		t, perr := time.Parse(time.RFC3339, req.PresaleStartDate)
		if perr == nil {
			event.PresaleStartDate = &t
		}
	}
	if req.GeneralSaleDate != "" {
		t, perr := time.Parse(time.RFC3339, req.GeneralSaleDate)
		if perr == nil {
			event.GeneralSaleDate = &t
		}
	}

	return event, nil
}

func (h *EventController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	var req updateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := buildEventUpdates(&req)
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	event, err := h.eventService.Update(uint(id), updates)
	if err != nil {
		if isNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toEventResponse(event))
}

func (h *EventController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	if err := h.eventService.Delete(uint(id)); err != nil {
		if isNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func buildEventUpdates(req *updateEventRequest) map[string]interface{} {
	updates := make(map[string]interface{})

	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.ImageURL != nil {
		updates["image_url"] = *req.ImageURL
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.EventDate != "" {
		t, err := time.Parse(time.RFC3339, req.EventDate)
		if err == nil {
			updates["event_date"] = t
		}
	}
	if req.DurationMinutes != nil {
		updates["duration_minutes"] = *req.DurationMinutes
	}
	if req.Capacity != nil {
		updates["capacity"] = *req.Capacity
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.PresaleCode != nil {
		updates["presale_code"] = *req.PresaleCode
	}
	if req.PresaleStartDate != nil {
		t, err := time.Parse(time.RFC3339, *req.PresaleStartDate)
		if err == nil {
			updates["presale_start_date"] = t
		}
	}
	if req.GeneralSaleDate != nil {
		t, err := time.Parse(time.RFC3339, *req.GeneralSaleDate)
		if err == nil {
			updates["general_sale_date"] = t
		}
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	return updates
}
