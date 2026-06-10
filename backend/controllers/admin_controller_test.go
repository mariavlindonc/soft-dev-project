package controllers

import (
	"encoding/json"
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
}
