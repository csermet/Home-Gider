package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/caner/home-gider/internal/models"
	"github.com/caner/home-gider/internal/services"
	"github.com/labstack/echo/v4"
)

type ExpenseHandler struct {
	service *services.ExpenseService
}

func NewExpenseHandler(service *services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: service}
}

type CreateExpenseRequest struct {
	CategoryID    uint    `json:"category_id"`
	Description   string  `json:"description"`
	Amount        float64 `json:"amount"`
	ExpenseDate   string  `json:"expense_date"`
	IsShared      bool    `json:"is_shared"`
	SplitRatio    float64 `json:"split_ratio"`
}

func (h *ExpenseHandler) List(c echo.Context) error {
	now := time.Now()
	month, _ := strconv.Atoi(c.QueryParam("month"))
	year, _ := strconv.Atoi(c.QueryParam("year"))
	if month == 0 {
		month = int(now.Month())
	}
	if year == 0 {
		year = now.Year()
	}

	expenses, err := h.service.List(month, year)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Giderler yüklenemedi"})
	}
	return c.JSON(http.StatusOK, expenses)
}

func (h *ExpenseHandler) Create(c echo.Context) error {
	var req CreateExpenseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
	}

	userID := c.Get("user_id").(uint)

	date, err := time.Parse("2006-01-02", req.ExpenseDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz tarih formatı (YYYY-MM-DD)"})
	}

	if req.SplitRatio == 0 {
		req.SplitRatio = 50
	}

	expense := &models.Expense{
		CreatedBy:   userID,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		Amount:      req.Amount,
		ExpenseDate: date,
		IsShared:    req.IsShared,
		SplitRatio:  req.SplitRatio,
	}

	if err := h.service.Create(expense); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gider oluşturulamadı"})
	}

	created, _ := h.service.GetByID(expense.ID)
	return c.JSON(http.StatusCreated, created)
}

func (h *ExpenseHandler) Update(c echo.Context) error {
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

	updated, _ := h.service.GetByID(uint(id))
	return c.JSON(http.StatusOK, updated)
}

func (h *ExpenseHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz ID"})
	}

	userID := c.Get("user_id").(uint)
	if err := h.service.Delete(uint(id), userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Gider silindi"})
}

func (h *ExpenseHandler) Approve(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz ID"})
	}

	userID := c.Get("user_id").(uint)
	isAdmin, _ := c.Get("is_admin").(bool)
	if err := h.service.Approve(uint(id), userID, isAdmin); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	updated, _ := h.service.GetByID(uint(id))
	return c.JSON(http.StatusOK, updated)
}

func (h *ExpenseHandler) Reject(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz ID"})
	}

	userID := c.Get("user_id").(uint)
	isAdmin, _ := c.Get("is_admin").(bool)
	if err := h.service.Reject(uint(id), userID, isAdmin); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	updated, _ := h.service.GetByID(uint(id))
	return c.JSON(http.StatusOK, updated)
}
