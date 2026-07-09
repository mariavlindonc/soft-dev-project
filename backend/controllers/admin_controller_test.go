package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGlobalReport(t *testing.T) {
	mockSvc := new(MockReportService)
	ctrl := NewAdminController(mockSvc)

	mockSvc.On("GetGlobalReport").Return(&services.GlobalReport{
		TotalEvents:      3,
		TotalTicketsSold: 150,
		EventReports: []services.EventReport{
			{EventID: 1, EventTitle: "E1", TotalCapacity: 100, TicketsSold: 80, Occupancy: 80},
		},
	}, nil)

	r := setupRouter()
	r.GET("/admin/reports", setAuthContext("admin", 1), ctrl.GetReports)

	req := httptest.NewRequest(http.MethodGet, "/admin/reports", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var resp globalReportResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, 3, resp.TotalEvents)
	assert.Equal(t, 150, resp.TotalTicketsSold)
}

func TestGetEventReport(t *testing.T) {
	t.Run("existing event returns report", func(t *testing.T) {
		mockSvc := new(MockReportService)
		ctrl := NewAdminController(mockSvc)

		mockSvc.On("GetEventReport", uint(1)).Return(&services.EventReport{
			EventID: 1, EventTitle: "Concert", TotalCapacity: 100,
			TicketsSold: 75, Occupancy: 75,
			Buyers: []services.BuyerInfo{
				{UserID: 10, Name: "Alice", Email: "a@t.com"},
			},
		}, nil)

		r := setupRouter()
		r.GET("/admin/reports/events/:id", setAuthContext("admin", 1), ctrl.GetEventReport)

		req := httptest.NewRequest(http.MethodGet, "/admin/reports/events/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		var resp eventReportResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, uint(1), resp.EventID)
		assert.Len(t, resp.Buyers, 1)
	})

	t.Run("not found returns 404", func(t *testing.T) {
		mockSvc := new(MockReportService)
		ctrl := NewAdminController(mockSvc)

		mockSvc.On("GetEventReport", uint(99)).Return(nil, services.ErrNotFound)

		r := setupRouter()
		r.GET("/admin/reports/events/:id", setAuthContext("admin", 1), ctrl.GetEventReport)

		req := httptest.NewRequest(http.MethodGet, "/admin/reports/events/99", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("invalid event id returns 400", func(t *testing.T) {
		ctrl := NewAdminController(new(MockReportService))
		r := setupRouter()
		r.GET("/admin/reports/events/:id", setAuthContext("admin", 1), ctrl.GetEventReport)

		req := httptest.NewRequest(http.MethodGet, "/admin/reports/events/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal error returns 500", func(t *testing.T) {
		mockSvc := new(MockReportService)
		ctrl := NewAdminController(mockSvc)

		mockSvc.On("GetEventReport", uint(1)).Return(nil, fmt.Errorf("db error"))

		r := setupRouter()
		r.GET("/admin/reports/events/:id", setAuthContext("admin", 1), ctrl.GetEventReport)

		req := httptest.NewRequest(http.MethodGet, "/admin/reports/events/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetGlobalReportErrors(t *testing.T) {
	t.Run("service error returns 500", func(t *testing.T) {
		mockSvc := new(MockReportService)
		ctrl := NewAdminController(mockSvc)

		mockSvc.On("GetGlobalReport").Return(nil, fmt.Errorf("db error"))

		r := setupRouter()
		r.GET("/admin/reports", setAuthContext("admin", 1), ctrl.GetReports)

		req := httptest.NewRequest(http.MethodGet, "/admin/reports", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGlobalReportResponseMapping(t *testing.T) {
	report := &services.GlobalReport{
		TotalEvents:      2,
		TotalTicketsSold: 50,
		EventReports: []services.EventReport{
			{EventID: 1, EventTitle: "E1", TotalCapacity: 100, TicketsSold: 30, Occupancy: 30.0},
			{EventID: 2, EventTitle: "E2", TotalCapacity: 200, TicketsSold: 20, Occupancy: 10.0},
		},
	}

	resp := toGlobalReportResponse(report)
	assert.Equal(t, 2, resp.TotalEvents)
	assert.Equal(t, 50, resp.TotalTicketsSold)
	assert.Len(t, resp.Events, 2)
	assert.Equal(t, "E1", resp.Events[0].Title)
	assert.Equal(t, 30.0, resp.Events[0].Occupancy)
}

func TestEventReportResponseMapping(t *testing.T) {
	report := &services.EventReport{
		EventID:       1,
		EventTitle:    "Concert",
		TotalCapacity: 100,
		TicketsSold:   75,
		Occupancy:     75.0,
		Buyers: []services.BuyerInfo{
			{UserID: 10, Name: "Alice", Email: "a@t.com"},
			{UserID: 20, Name: "Bob", Email: "b@t.com"},
		},
	}

	resp := toEventReportResponse(report)
	assert.Equal(t, uint(1), resp.EventID)
	assert.Equal(t, "Concert", resp.Title)
	assert.Equal(t, 75.0, resp.Occupancy)
	assert.Len(t, resp.Buyers, 2)
	assert.Equal(t, "Alice", resp.Buyers[0].Name)
}
