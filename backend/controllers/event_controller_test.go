package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"backend/domain"
	"backend/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetAllEvents(t *testing.T) {
	mockSvc := new(MockEventService)
	ctrl := NewEventController(mockSvc)

	mockSvc.On("GetAll", mock.Anything).Return([]domain.Event{
		{ID: 1, Title: "Concert", EventDate: time.Now()},
	}, nil)

	r := setupRouter()
	r.GET("/events", ctrl.GetAll)

	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp []eventResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "Concert", resp[0].Title)
}

func TestGetEventByID(t *testing.T) {
	t.Run("existing event returns 200", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("GetByID", uint(1)).Return(&domain.Event{
			ID: 1, Title: "Concert", EventDate: time.Now(),
		}, nil)

		r := setupRouter()
		r.GET("/events/:id", ctrl.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/events/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not found returns 404", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("GetByID", uint(99)).Return(nil, services.ErrNotFound)

		r := setupRouter()
		r.GET("/events/:id", ctrl.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/events/99", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		ctrl := NewEventController(new(MockEventService))
		r := setupRouter()
		r.GET("/events/:id", ctrl.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/events/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCreateEvent(t *testing.T) {
	t.Run("valid event returns 201", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("Create", mock.Anything).Return(&domain.Event{
			ID: 1, Title: "New Event", EventDate: time.Now(),
		}, nil)

		r := setupRouter()
		r.POST("/events", ctrl.Create)

		future := time.Now().Add(48 * time.Hour).Format(time.RFC3339)
		body := `{"title":"New Event","event_date":"` + future + `","capacity":100,"price":50}`
		req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("missing title returns 400", func(t *testing.T) {
		ctrl := NewEventController(new(MockEventService))
		r := setupRouter()
		r.POST("/events", ctrl.Create)

		body := `{"capacity":100}`
		req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetSaleStatus(t *testing.T) {
	mockSvc := new(MockEventService)
	ctrl := NewEventController(mockSvc)

	mockSvc.On("GetByID", uint(1)).Return(&domain.Event{
		ID: 1, Title: "Concert",
		PresaleActive:    true,
		PresaleStartDate: timePtr(time.Now().Add(-2 * time.Hour)),
		GeneralSaleDate:  timePtr(time.Now().Add(2 * time.Hour)),
		EventDate:        time.Now().Add(24 * time.Hour),
		Status:           "presale",
	}, nil)

	r := setupRouter()
	r.GET("/events/:id/sale-status", ctrl.GetSaleStatus)

	req := httptest.NewRequest(http.MethodGet, "/events/1/sale-status", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp saleStatusResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, domain.SalePhase("presale"), resp.Phase)
}

func timePtr(t time.Time) *time.Time { return &t }

func TestUpdateEvent(t *testing.T) {
	t.Run("valid update returns 200", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("Update", uint(1), mock.Anything).Return(&domain.Event{
			ID: 1, Title: "Updated", EventDate: time.Now(),
		}, nil)

		r := setupRouter()
		r.PUT("/events/:id", ctrl.Update)

		body := `{"title":"Updated"}`
		req := httptest.NewRequest(http.MethodPut, "/events/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not found returns 404", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("Update", uint(99), mock.Anything).Return(nil, services.ErrNotFound)

		r := setupRouter()
		r.PUT("/events/:id", ctrl.Update)

		body := `{"title":"Nope"}`
		req := httptest.NewRequest(http.MethodPut, "/events/99", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDeleteEvent(t *testing.T) {
	t.Run("cancel event returns 204", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("Cancel", uint(1)).Return(nil)

		r := setupRouter()
		r.DELETE("/events/:id", ctrl.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/events/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("not found returns 404", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("Cancel", uint(99)).Return(services.ErrNotFound)

		r := setupRouter()
		r.DELETE("/events/:id", ctrl.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/events/99", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("invalid id returns 400", func(t *testing.T) {
		ctrl := NewEventController(new(MockEventService))
		r := setupRouter()
		r.DELETE("/events/:id", ctrl.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/events/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetAllEventsFilters(t *testing.T) {
	t.Run("with category filter", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("GetAll", mock.MatchedBy(func(f domain.EventFilters) bool {
			return f.Category == "Rock"
		})).Return([]domain.Event{
			{ID: 1, Title: "Rock Concert", EventDate: time.Now()},
		}, nil)

		r := setupRouter()
		r.GET("/events", ctrl.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/events?category=Rock", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("with date filter", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("GetAll", mock.MatchedBy(func(f domain.EventFilters) bool {
			return f.DateFrom != nil && f.DateTo != nil
		})).Return([]domain.Event{}, nil)

		r := setupRouter()
		r.GET("/events", ctrl.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/events?date=2026-06-15", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("with date_from and date_to filters", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("GetAll", mock.MatchedBy(func(f domain.EventFilters) bool {
			return f.DateFrom != nil && f.DateTo != nil
		})).Return([]domain.Event{}, nil)

		r := setupRouter()
		r.GET("/events", ctrl.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/events?date_from=2026-06-01&date_to=2026-06-30", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("with min_price and max_price filters", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("GetAll", mock.MatchedBy(func(f domain.EventFilters) bool {
			return f.MinPrice != nil && f.MaxPrice != nil
		})).Return([]domain.Event{}, nil)

		r := setupRouter()
		r.GET("/events", ctrl.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/events?min_price=10&max_price=100", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("with invalid date filter is ignored", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("GetAll", mock.MatchedBy(func(f domain.EventFilters) bool {
			return f.DateFrom == nil && f.DateTo == nil
		})).Return([]domain.Event{}, nil)

		r := setupRouter()
		r.GET("/events", ctrl.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/events?date=invalid-date", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("with invalid min_price is ignored", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("GetAll", mock.MatchedBy(func(f domain.EventFilters) bool {
			return f.MinPrice == nil
		})).Return([]domain.Event{}, nil)

		r := setupRouter()
		r.GET("/events", ctrl.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/events?min_price=abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("with invalid max_price is ignored", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("GetAll", mock.MatchedBy(func(f domain.EventFilters) bool {
			return f.MaxPrice == nil
		})).Return([]domain.Event{}, nil)

		r := setupRouter()
		r.GET("/events", ctrl.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/events?max_price=abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("service error returns 500", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("GetAll", mock.Anything).Return([]domain.Event{}, fmt.Errorf("db error"))

		r := setupRouter()
		r.GET("/events", ctrl.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/events", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("invalid date_from is ignored", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("GetAll", mock.MatchedBy(func(f domain.EventFilters) bool {
			return f.DateFrom == nil
		})).Return([]domain.Event{}, nil)

		r := setupRouter()
		r.GET("/events", ctrl.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/events?date_from=bad", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid date_to is ignored", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("GetAll", mock.MatchedBy(func(f domain.EventFilters) bool {
			return f.DateTo == nil
		})).Return([]domain.Event{}, nil)

		r := setupRouter()
		r.GET("/events", ctrl.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/events?date_to=bad", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestCreateEventErrors(t *testing.T) {
	t.Run("invalid event_date format returns 400", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		r := setupRouter()
		r.POST("/events", ctrl.Create)

		body := `{"title":"Event","event_date":"bad-date","capacity":100,"price":50}`
		req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service returns ErrInvalidInput returns 400", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("Create", mock.Anything).Return(nil, services.ErrInvalidInput)

		r := setupRouter()
		r.POST("/events", ctrl.Create)

		future := time.Now().Add(48 * time.Hour).Format(time.RFC3339)
		body := `{"title":"Event","event_date":"` + future + `","capacity":100,"price":50}`
		req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service returns other error returns 500", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("Create", mock.Anything).Return(nil, fmt.Errorf("db error"))

		r := setupRouter()
		r.POST("/events", ctrl.Create)

		future := time.Now().Add(48 * time.Hour).Format(time.RFC3339)
		body := `{"title":"Event","event_date":"` + future + `","capacity":100,"price":50}`
		req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUpdateEventErrors(t *testing.T) {
	t.Run("invalid id returns 400", func(t *testing.T) {
		ctrl := NewEventController(new(MockEventService))
		r := setupRouter()
		r.PUT("/events/:id", ctrl.Update)

		body := `{"title":"Updated"}`
		req := httptest.NewRequest(http.MethodPut, "/events/abc", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid json returns 400", func(t *testing.T) {
		ctrl := NewEventController(new(MockEventService))
		r := setupRouter()
		r.PUT("/events/:id", ctrl.Update)

		req := httptest.NewRequest(http.MethodPut, "/events/1", strings.NewReader("not json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("cancelled event returns 400", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("Update", uint(1), mock.Anything).Return(nil, services.ErrEventCancelled)

		r := setupRouter()
		r.PUT("/events/:id", ctrl.Update)

		body := `{"title":"Updated"}`
		req := httptest.NewRequest(http.MethodPut, "/events/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal error returns 500", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("Update", uint(1), mock.Anything).Return(nil, fmt.Errorf("db error"))

		r := setupRouter()
		r.PUT("/events/:id", ctrl.Update)

		body := `{"title":"Updated"}`
		req := httptest.NewRequest(http.MethodPut, "/events/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetSaleStatusVariations(t *testing.T) {
	t.Run("invalid id returns 400", func(t *testing.T) {
		ctrl := NewEventController(new(MockEventService))
		r := setupRouter()
		r.GET("/events/:id/sale-status", ctrl.GetSaleStatus)

		req := httptest.NewRequest(http.MethodGet, "/events/abc/sale-status", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found returns 404", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("GetByID", uint(99)).Return(nil, services.ErrNotFound)

		r := setupRouter()
		r.GET("/events/:id/sale-status", ctrl.GetSaleStatus)

		req := httptest.NewRequest(http.MethodGet, "/events/99/sale-status", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("cancelled event returns cancelled message", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("GetByID", uint(1)).Return(&domain.Event{
			ID: 1, Status: "cancelled",
		}, nil)

		r := setupRouter()
		r.GET("/events/:id/sale-status", ctrl.GetSaleStatus)

		req := httptest.NewRequest(http.MethodGet, "/events/1/sale-status", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp saleStatusResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Evento cancelado", resp.Message)
	})

	t.Run("sold out event returns sold out message", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("GetByID", uint(1)).Return(&domain.Event{
			ID: 1, Status: "active", Capacity: 100, TicketsSold: 100,
		}, nil)

		r := setupRouter()
		r.GET("/events/:id/sale-status", ctrl.GetSaleStatus)

		req := httptest.NewRequest(http.MethodGet, "/events/1/sale-status", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp saleStatusResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Entradas agotadas", resp.Message)
	})

	t.Run("no presale event returns empty message", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("GetByID", uint(1)).Return(&domain.Event{
			ID: 1, Status: "active", Capacity: 100,
		}, nil)

		r := setupRouter()
		r.GET("/events/:id/sale-status", ctrl.GetSaleStatus)

		req := httptest.NewRequest(http.MethodGet, "/events/1/sale-status", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp saleStatusResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, domain.PhaseNoPresale, resp.Phase)
	})

	t.Run("service error returns 500", func(t *testing.T) {
		mockSvc := new(MockEventService)
		ctrl := NewEventController(mockSvc)

		mockSvc.On("GetByID", uint(1)).Return(nil, fmt.Errorf("db error"))

		r := setupRouter()
		r.GET("/events/:id/sale-status", ctrl.GetSaleStatus)

		req := httptest.NewRequest(http.MethodGet, "/events/1/sale-status", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestSalePhaseMessage(t *testing.T) {
	now := time.Now()

	t.Run("NotYetOpen with presale start", func(t *testing.T) {
		msg := salePhaseMessage(domain.PhaseNotYetOpen, &now, nil)
		assert.Contains(t, msg, "Las ventas aún no abrieron")
	})

	t.Run("NotYetOpen without presale start", func(t *testing.T) {
		msg := salePhaseMessage(domain.PhaseNotYetOpen, nil, nil)
		assert.Contains(t, msg, "Las ventas aún no han comenzado")
	})

	t.Run("Presale with general sale", func(t *testing.T) {
		msg := salePhaseMessage(domain.PhasePresale, nil, &now)
		assert.Contains(t, msg, "Pre-venta activa")
	})

	t.Run("Presale without general sale", func(t *testing.T) {
		msg := salePhaseMessage(domain.PhasePresale, nil, nil)
		assert.Contains(t, msg, "Preventa activa")
	})

	t.Run("Public with general sale", func(t *testing.T) {
		msg := salePhaseMessage(domain.PhasePublic, nil, &now)
		assert.Contains(t, msg, "Venta general disponible")
	})

	t.Run("Public without general sale", func(t *testing.T) {
		msg := salePhaseMessage(domain.PhasePublic, nil, nil)
		assert.Contains(t, msg, "Venta general abierta")
	})

	t.Run("NoPresale returns empty", func(t *testing.T) {
		msg := salePhaseMessage(domain.PhaseNoPresale, nil, nil)
		assert.Equal(t, "", msg)
	})

	t.Run("unknown phase returns empty", func(t *testing.T) {
		msg := salePhaseMessage(domain.SalePhase("unknown"), nil, nil)
		assert.Equal(t, "", msg)
	})
}

func TestGetByIDInternalError(t *testing.T) {
	mockSvc := new(MockEventService)
	ctrl := NewEventController(mockSvc)

	mockSvc.On("GetByID", uint(1)).Return(nil, fmt.Errorf("db error"))

	r := setupRouter()
	r.GET("/events/:id", ctrl.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/events/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
