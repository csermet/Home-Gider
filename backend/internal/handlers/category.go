package handlers

import (
	"net/http"

	"github.com/caner/home-gider/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	db *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

func (h *CategoryHandler) List(c echo.Context) error {
	var categories []models.Category
	h.db.Order("id ASC").Find(&categories)
	return c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) Create(c echo.Context) error {
	var cat models.Category
	if err := c.Bind(&cat); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
	}
	if err := h.db.Create(&cat).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Kategori oluşturulamadı"})
	}
	return c.JSON(http.StatusCreated, cat)
}
