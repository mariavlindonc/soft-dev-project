package controllers

import (
	"net/http"
	"strconv"
	"time"

	"backend/domain"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	ticketService services.TicketService
}

func NewTicketController(s services.TicketService) *TicketController {
	return &TicketController{ticketService: s}
}

type purchaseTicketRequest struct {
	EventID uint `json:"event_id" binding:"required"`
}

type transferTicketRequest struct {
	ToUserID uint `json:"to_user_id" binding:"required"`
}

type ticketResponse struct {
	ID            uint    `json:"id"`
	UserID        uint    `json:"user_id"`
	EventID       uint    `json:"event_id"`
	Status        string  `json:"status"`
	PurchasePrice float64 `json:"purchase_price"`
	PurchasedAt   string  `json:"purchased_at"`
}

func toTicketResponse(t *domain.Ticket) ticketResponse {
	return ticketResponse{
		ID:            t.ID,
		UserID:        t.UserID,
		EventID:       t.EventID,
		Status:        t.Status,
		PurchasePrice: t.PurchasePrice,
		PurchasedAt:   t.PurchasedAt.Format(time.RFC3339),
	}
}

func (h *TicketController) Purchase(c *gin.Context) {
	var req purchaseTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")

	ticket, err := h.ticketService.Purchase(userID, req.EventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toTicketResponse(ticket))
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

	if err := h.ticketService.Transfer(uint(id), userID, req.ToUserID); err != nil {
		if isNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket transferred successfully"})
}
