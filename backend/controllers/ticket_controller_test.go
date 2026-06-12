package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"backend/domain"
	"backend/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setAuthContext(role string, userID uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}

func TestPurchaseTicket(t *testing.T) {
	t.Run("valid purchase returns 201", func(t *testing.T) {
		mockSvc := new(MockTicketService)
		ctrl := NewTicketController(mockSvc)

		mockSvc.On("Purchase", uint(10), services.PurchaseInput{
			EventID: 1, Quantity: 2, PresaleCode: "",
		}).Return([]*domain.Ticket{
			{ID: 1, UserID: 10, EventID: 1, Status: "active", PurchasePrice: 50, PurchasedAt: time.Now()},
			{ID: 2, UserID: 10, EventID: 1, Status: "active", PurchasePrice: 50, PurchasedAt: time.Now()},
		}, nil)

		r := setupRouter()
		r.POST("/tickets/purchase", setAuthContext("client", 10), ctrl.Purchase)

		body := `{"event_id":1,"quantity":2}`
		req := httptest.NewRequest(http.MethodPost, "/tickets/purchase", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
		var resp []ticketResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Len(t, resp, 2)
		assert.Equal(t, uint(1), resp[0].ID)
		assert.Equal(t, "active", resp[0].Status)
	})

	t.Run("cancelled event returns 400", func(t *testing.T) {
		mockSvc := new(MockTicketService)
		ctrl := NewTicketController(mockSvc)

		mockSvc.On("Purchase", uint(10), mock.Anything).Return(nil, services.ErrEventCancelled)

		r := setupRouter()
		r.POST("/tickets/purchase", setAuthContext("client", 10), ctrl.Purchase)

		body := `{"event_id":1,"quantity":1}`
		req := httptest.NewRequest(http.MethodPost, "/tickets/purchase", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing event_id returns 400", func(t *testing.T) {
		ctrl := NewTicketController(new(MockTicketService))
		r := setupRouter()
		r.POST("/tickets/purchase", setAuthContext("client", 10), ctrl.Purchase)

		body := `{}`
		req := httptest.NewRequest(http.MethodPost, "/tickets/purchase", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetMyTickets(t *testing.T) {
	mockSvc := new(MockTicketService)
	ctrl := NewTicketController(mockSvc)

	mockSvc.On("GetByUser", uint(10)).Return([]domain.Ticket{
		{ID: 1, UserID: 10, EventID: 1, Status: "active", PurchasedAt: time.Now()},
	}, nil)

	r := setupRouter()
	r.GET("/tickets", setAuthContext("client", 10), ctrl.GetMyTickets)

	req := httptest.NewRequest(http.MethodGet, "/tickets", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp []ticketResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 1)
}

func TestCancelTicket(t *testing.T) {
	t.Run("cancel own ticket returns 204", func(t *testing.T) {
		mockSvc := new(MockTicketService)
		ctrl := NewTicketController(mockSvc)

		mockSvc.On("Cancel", uint(1), uint(10)).Return(nil)

		r := setupRouter()
		r.PATCH("/tickets/:id/cancel", setAuthContext("client", 10), ctrl.Cancel)

		req := httptest.NewRequest(http.MethodPatch, "/tickets/1/cancel", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("cancel not found ticket returns 404", func(t *testing.T) {
		mockSvc := new(MockTicketService)
		ctrl := NewTicketController(mockSvc)

		mockSvc.On("Cancel", uint(99), uint(10)).Return(services.ErrNotFound)

		r := setupRouter()
		r.PATCH("/tickets/:id/cancel", setAuthContext("client", 10), ctrl.Cancel)

		req := httptest.NewRequest(http.MethodPatch, "/tickets/99/cancel", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestTransferTicket(t *testing.T) {
	t.Run("valid transfer returns 200", func(t *testing.T) {
		mockSvc := new(MockTicketService)
		ctrl := NewTicketController(mockSvc)

		mockSvc.On("Transfer", uint(1), uint(10), services.TransferInput{
			ToUserEmail: "other@test.com",
		}).Return(nil)

		r := setupRouter()
		r.PATCH("/tickets/:id/transfer", setAuthContext("client", 10), ctrl.Transfer)

		body := `{"to_user_email":"other@test.com"}`
		req := httptest.NewRequest(http.MethodPatch, "/tickets/1/transfer", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("missing email returns 400", func(t *testing.T) {
		ctrl := NewTicketController(new(MockTicketService))
		r := setupRouter()
		r.PATCH("/tickets/:id/transfer", setAuthContext("client", 10), ctrl.Transfer)

		body := `{}`
		req := httptest.NewRequest(http.MethodPatch, "/tickets/1/transfer", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
