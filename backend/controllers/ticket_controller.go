package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"backend/domain"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	ticketService services.TicketServiceInterface
}

func NewTicketController(s services.TicketServiceInterface) *TicketController {
	return &TicketController{ticketService: s}
}

type purchaseTicketRequest struct {
	EventID     uint   `json:"event_id"     binding:"required"`
	Quantity    uint   `json:"quantity"     binding:"required,min=1"`
	PresaleCode string `json:"presale_code"`
}

type transferTicketRequest struct {
	ToUserEmail string `json:"to_user_email" binding:"required,email"`
}

type ticketResponse struct {
	ID            uint    `json:"id"`
	UserID        uint    `json:"user_id"`
	EventID       uint    `json:"event_id"`
	Status        string  `json:"status"`
	PurchasePrice float64 `json:"purchase_price"`
	PurchasedAt   string  `json:"purchased_at"`
	EventTitle    string  `json:"event_title"`
	EventDate     string  `json:"event_date"`
}

func toTicketResponse(t *domain.Ticket) ticketResponse {
	eventTitle := ""
	eventDate := ""
	if t.Event.Title != "" {
		eventTitle = t.Event.Title
		eventDate = t.Event.EventDate.Format(time.RFC3339)
	}

	return ticketResponse{
		ID:            t.ID,
		UserID:        t.UserID,
		EventID:       t.EventID,
		Status:        t.Status,
		PurchasePrice: t.PurchasePrice,
		PurchasedAt:   t.PurchasedAt.Format(time.RFC3339),
		EventTitle:    eventTitle,
		EventDate:     eventDate,
	}
}

func (h *TicketController) Purchase(c *gin.Context) {
	var req purchaseTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")

	tickets, err := h.ticketService.Purchase(userID, services.PurchaseInput{EventID: req.EventID, Quantity: req.Quantity, PresaleCode: req.PresaleCode})
	if err != nil {
		if isNotFound(err) || errors.Is(err, services.ErrEventCancelled) || errors.Is(err, services.ErrNoCapacity) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]ticketResponse, len(tickets))
	for i, t := range tickets {
		resp[i] = toTicketResponse(t)
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *TicketController) GetMyTickets(c *gin.Context) {
	userID := c.GetUint("userID")

	tickets, err := h.ticketService.GetByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]ticketResponse, 0, len(tickets))
	for _, t := range tickets {
		resp = append(resp, toTicketResponse(&t))
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TicketController) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	userID := c.GetUint("userID")

	if err := h.ticketService.Cancel(uint(id), userID); err != nil {
		if isNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *TicketController) Transfer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	var req transferTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")

	if err := h.ticketService.Transfer(uint(id), userID, services.TransferInput{ToUserEmail: req.ToUserEmail}); err != nil {
		if isNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket transferred successfully"})
}
