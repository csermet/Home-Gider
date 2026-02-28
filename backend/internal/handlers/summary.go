package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/caner/home-gider/internal/services"
	"github.com/labstack/echo/v4"
)

type SummaryHandler struct {
	service *services.SettlementService
}

func NewSummaryHandler(service *services.SettlementService) *SummaryHandler {
	return &SummaryHandler{service: service}
}

func (h *SummaryHandler) GetSummary(c echo.Context) error {
	now := time.Now()
	month, _ := strconv.Atoi(c.QueryParam("month"))
	year, _ := strconv.Atoi(c.QueryParam("year"))
	if month == 0 {
		month = int(now.Month())
	}
	if year == 0 {
		year = now.Year()
	}

	sharedOnly := c.QueryParam("shared_only") == "true"

	summary, err := h.service.GetMonthlySummary(month, year, sharedOnly)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Özet hesaplanamadı"})
	}
	return c.JSON(http.StatusOK, summary)
}
