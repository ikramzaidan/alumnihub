package handler

import (
	"errors"
	"net/http"

	"alumnihub/internal/service"
	"alumnihub/internal/utils"
)

func NewDashboardHandler(dashboardService *service.DashboardService) *Handler {
	return &Handler{DashboardService: dashboardService}
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	if h.DashboardService == nil {
		_ = utils.ErrorJSON(w, errors.New("dashboard service not initialized"))
		return
	}

	payload, err := h.DashboardService.GetSummary()
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, payload)
}
