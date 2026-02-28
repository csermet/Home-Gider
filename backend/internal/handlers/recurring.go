package handlers

import (
	"net/http"
	"strconv"

	"github.com/caner/home-gider/internal/models"
	"github.com/caner/home-gider/internal/services"
	"github.com/labstack/echo/v4"
)

type RecurringHandler struct {
	service *services.RecurringService
}

func NewRecurringHandler(service *services.RecurringService) *RecurringHandler {
	return &RecurringHandler{service: service}
}

type CreateRecurringRequest struct {
	CategoryID       uint    `json:"category_id"`
	Description      string  `json:"description"`
	Amount           float64 `json:"amount"`
	TotalAmount      float64 `json:"total_amount"`
	Type             string  `json:"type"`
	InstallmentCount int     `json:"installment_count"`
	IsShared         bool    `json:"is_shared"`
	SplitRatio       float64 `json:"split_ratio"`
}

func (h *RecurringHandler) List(c echo.Context) error {
	items, err := h.service.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Şablonlar yüklenemedi"})
	}
	return c.JSON(http.StatusOK, items)
}

func (h *RecurringHandler) Create(c echo.Context) error {
	var req CreateRecurringRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
	}

	userID := c.Get("user_id").(uint)

	if req.SplitRatio == 0 {
		req.SplitRatio = 50
	}

	item := &models.RecurringExpense{
		CreatedBy:   userID,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		Amount:      req.Amount,
		Type:        models.RecurringType(req.Type),
		IsShared:    req.IsShared,
		SplitRatio:  req.SplitRatio,
	}

	if req.TotalAmount > 0 {
		item.TotalAmount = &req.TotalAmount
	}
	if req.InstallmentCount > 0 {
		item.InstallmentCount = &req.InstallmentCount
	}

	if err := h.service.Create(item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, item)
}

func (h *RecurringHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz ID"})
	}

	var updates map[string]interface{}
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
	}

	userID := c.Get("user_id").(uint)
	if err := h.service.Update(uint(id), userID, updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Güncellendi"})
}

func (h *RecurringHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz ID"})
	}

	userID := c.Get("user_id").(uint)
	if err := h.service.Delete(uint(id), userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Deaktif edildi"})
}

func (h *RecurringHandler) Approve(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz ID"})
	}

	userID := c.Get("user_id").(uint)
	isAdmin, _ := c.Get("is_admin").(bool)
	if err := h.service.Approve(uint(id), userID, isAdmin); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Onaylandı"})
}

func (h *RecurringHandler) Reject(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz ID"})
	}

	userID := c.Get("user_id").(uint)
	isAdmin, _ := c.Get("is_admin").(bool)
	if err := h.service.Reject(uint(id), userID, isAdmin); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Reddedildi"})
}
