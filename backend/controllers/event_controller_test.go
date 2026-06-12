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
}
