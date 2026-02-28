package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/caner/home-gider/internal/services"
	"github.com/labstack/echo/v4"
)

type SettlementHandler struct {
	service *services.SettlementService
}

func NewSettlementHandler(service *services.SettlementService) *SettlementHandler {
	return &SettlementHandler{service: service}
}

func (h *SettlementHandler) ListPayments(c echo.Context) error {
	now := time.Now()
	month, _ := strconv.Atoi(c.QueryParam("month"))
	year, _ := strconv.Atoi(c.QueryParam("year"))
	if month == 0 {
		month = int(now.Month())
	}
	if year == 0 {
		year = now.Year()
	}

	payments, err := h.service.GetPayments(month, year)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ödemeler yüklenemedi"})
	}
	return c.JSON(http.StatusOK, payments)
}

type AddPaymentRequest struct {
	Month   int     `json:"month"`
	Year    int     `json:"year"`
	PayerID uint    `json:"payer_id"`
	PayeeID uint    `json:"payee_id"`
	Amount  float64 `json:"amount"`
}

func (h *SettlementHandler) AddPayment(c echo.Context) error {
	var req AddPaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
	}
	if req.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tutar sıfırdan büyük olmalı"})
	}

	// Sadece borçlu olan kişi ödeme ekleyebilir
	userID := c.Get("user_id").(uint)
	if req.PayerID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Sadece borçlu kişi ödeme ekleyebilir"})
	}

	if err := h.service.AddPayment(req.Month, req.Year, req.PayerID, req.PayeeID, req.Amount); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "Ödeme kaydedildi"})
}

func (h *SettlementHandler) DeletePayment(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz ID"})
	}
	if err := h.service.DeletePayment(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Ödeme silinemedi"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Ödeme silindi"})
}
