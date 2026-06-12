package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"backend/domain"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	eventService services.EventServiceInterface
}

func NewEventController(s services.EventServiceInterface) *EventController {
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
	PresaleActive    bool    `json:"presale_active"`
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
	PresaleActive    *bool    `json:"presale_active"`
	PresaleCode      *string  `json:"presale_code"`
	PresaleStartDate *string  `json:"presale_start_date"`
	GeneralSaleDate  *string  `json:"general_sale_date"`
	Status           *string  `json:"status"             binding:"omitempty,oneof=active presale sold_out cancelled"`
}

type eventResponse struct {
	ID               uint    `json:"id"`
	Title            string  `json:"title"`
	Description      *string `json:"description,omitempty"`
	ImageURL         *string `json:"image_url,omitempty"`
	Category         *string `json:"category,omitempty"`
	Location         *string `json:"location,omitempty"`
	EventDate        string  `json:"event_date"`
	DurationMinutes  int     `json:"duration_minutes"`
	Capacity         int     `json:"capacity"`
	TicketsSold      int     `json:"tickets_sold"`
	Price            float64 `json:"price"`
	Status           string  `json:"status"`
	PresaleActive    bool    `json:"presale_active"`
	PresaleStartDate *string `json:"presale_start_date,omitempty"`
	GeneralSaleDate  *string `json:"general_sale_date,omitempty"`
	CreatedByID      uint    `json:"created_by_id"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

func toEventResponse(e *domain.Event) eventResponse {
	return eventResponse{
		ID:               e.ID,
		Title:            e.Title,
		Description:      e.Description,
		ImageURL:         e.ImageURL,
		Category:         e.Category,
		Location:         e.Location,
		EventDate:        e.EventDate.Format(time.RFC3339),
		DurationMinutes:  e.DurationMinutes,
		Capacity:         e.Capacity,
		TicketsSold:      e.TicketsSold,
		Price:            e.Price,
		Status:           e.Status,
		PresaleActive:    e.PresaleActive,
		PresaleStartDate: timePtrToString(e.PresaleStartDate),
		GeneralSaleDate:  timePtrToString(e.GeneralSaleDate),
		CreatedByID:      e.CreatedByID,
		CreatedAt:        e.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        e.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *EventController) GetAll(c *gin.Context) {
	filters := domain.EventFilters{
		Category: c.Query("category"),
	}
	if dateStr := c.Query("date"); dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			filters.DateFrom = &parsed
			end := parsed.Add(24 * time.Hour)
			filters.DateTo = &end
		}
	}
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		parsed, err := time.Parse("2006-01-02", dateFromStr)
		if err == nil {
			filters.DateFrom = &parsed
		}
	}
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		parsed, err := time.Parse("2006-01-02", dateToStr)
		if err == nil {
			filters.DateTo = &parsed
		}
	}
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filters.MinPrice = &minPrice
		}
	}
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filters.MaxPrice = &maxPrice
		}
	}

	events, err := h.eventService.GetAll(filters)
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

	eventDate, err := time.Parse(time.RFC3339, req.EventDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_date, use RFC3339"})
		return
	}

	var presaleStart, generalSale *time.Time
	if req.PresaleStartDate != "" {
		t, perr := time.Parse(time.RFC3339, req.PresaleStartDate)
		if perr == nil {
			presaleStart = &t
		}
	}
	if req.GeneralSaleDate != "" {
		t, perr := time.Parse(time.RFC3339, req.GeneralSaleDate)
		if perr == nil {
			generalSale = &t
		}
	}

	created, err := h.eventService.Create(services.CreateEventInput{
		Title:            req.Title,
		Date:             eventDate,
		Duration:         req.DurationMinutes,
		Capacity:         req.Capacity,
		Price:            req.Price,
		Category:         req.Category,
		Description:      req.Description,
		Location:         req.Location,
		ImageURL:         req.ImageURL,
		PresaleActive:    req.PresaleActive,
		PresaleCode:      req.PresaleCode,
		PresaleStartDate: presaleStart,
		GeneralSaleDate:  generalSale,
		CreatedByID:      c.GetUint("userID"),
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toEventResponse(created))
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

	input := services.UpdateEventInput{
		Title:            optionalString(req.Title),
		Description:      req.Description,
		ImageURL:         req.ImageURL,
		Category:         req.Category,
		Location:         req.Location,
		Duration:         req.DurationMinutes,
		Capacity:         req.Capacity,
		Price:            req.Price,
		PresaleActive:    req.PresaleActive,
		PresaleCode:      req.PresaleCode,
		PresaleStartDate: req.PresaleStartDate,
		GeneralSaleDate:  req.GeneralSaleDate,
		Status:           req.Status,
	}
	if req.EventDate != "" {
		t, perr := time.Parse(time.RFC3339, req.EventDate)
		if perr == nil {
			input.Date = &t
		}
	}
	input.PresaleStartDate = req.PresaleStartDate
	input.GeneralSaleDate = req.GeneralSaleDate

	event, err := h.eventService.Update(uint(id), input)
	if err != nil {
		if isNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}
		if errors.Is(err, services.ErrEventCancelled) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	if err := h.eventService.Cancel(uint(id)); err != nil {
		if isNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

type saleStatusResponse struct {
	Phase            domain.SalePhase `json:"phase"`
	PresaleStartDate *string          `json:"presale_start_date,omitempty"`
	GeneralSaleDate  *string          `json:"general_sale_date,omitempty"`
	Message          string           `json:"message"`
}

func salePhaseMessage(phase domain.SalePhase, presaleStart, generalSale *time.Time) string {
	switch phase {
	case domain.PhaseNotYetOpen:
		if presaleStart != nil {
			return fmt.Sprintf("Las ventas aún no abrieron. La preventa comienza el %s.", presaleStart.Format("02/01/2006 15:04"))
		}
		return "Las ventas aún no han comenzado."
	case domain.PhasePresale:
		if generalSale != nil {
			return fmt.Sprintf("Pre-venta activa. La venta general empieza en %s. Se require un código de acceso para comprar.", generalSale.Format("02/01/2006 15:04"))
		}
		return "Preventa activa. Se requiere un código de acceso para comprar."
	case domain.PhasePublic:
		if generalSale != nil {
			return fmt.Sprintf("Venta general disponible desde el %s. No se requiere código de acceso.", generalSale.Format("02/01/2006 15:04"))
		}
		return "Venta general abierta. No se requiere código de acceso."
	case domain.PhaseNoPresale:
		return ""
	default:
		return ""
	}
}

func (h *EventController) GetSaleStatus(c *gin.Context) {
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

	var phase domain.SalePhase
	var msg string
	if event.Status == "cancelled" {
		phase = domain.PhaseNotYetOpen
		msg = "Evento cancelado"
	} else if event.Capacity > 0 && event.TicketsSold >= event.Capacity {
		phase = domain.PhaseNotYetOpen
		msg = "Entradas agotadas"
	} else {
		phase = event.CurrentSalePhase(time.Now())
		msg = salePhaseMessage(phase, event.PresaleStartDate, event.GeneralSaleDate)
	}
	resp := saleStatusResponse{
		Phase:            phase,
		PresaleStartDate: timePtrToString(event.PresaleStartDate),
		GeneralSaleDate:  timePtrToString(event.GeneralSaleDate),
		Message:          msg,
	}

	c.JSON(http.StatusOK, resp)
}
