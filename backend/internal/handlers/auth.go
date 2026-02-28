package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/caner/home-gider/internal/services"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
	}

	user, token, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Kullanıcı adı veya şifre hatalı"})
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   30 * 24 * 60 * 60,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	}
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, map[string]string{"message": "Çıkış yapıldı"})
}

func (h *AuthHandler) Me(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	user, err := h.service.GetUser(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Kullanıcı bulunamadı"})
	}
	return c.JSON(http.StatusOK, user)
}

type ChangePasswordRequest struct {
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) ChangePassword(c echo.Context) error {
	var req ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
	}
	if len(req.NewPassword) < 4 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Şifre en az 4 karakter olmalı"})
	}

	userID := c.Get("user_id").(uint)
	if err := h.service.ChangePassword(userID, req.NewPassword); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Şifre değiştirilemedi"})
	}

	// Yeni token üret (must_change_password güncellendi)
	user, token, err := h.service.Login(c.Get("username").(string), req.NewPassword)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{"message": "Şifre değiştirildi, tekrar giriş yapın"})
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   30 * 24 * 60 * 60,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Şifre değiştirildi",
		"user":    user,
	})
}

// Admin: Tüm kullanıcıları listele
func (h *AuthHandler) ListUsers(c echo.Context) error {
	users, err := h.service.ListUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Kullanıcılar yüklenemedi"})
	}
	return c.JSON(http.StatusOK, users)
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password"`
}

// Admin: Kullanıcı şifresini sıfırla
func (h *AuthHandler) AdminResetPassword(c echo.Context) error {
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz ID"})
	}

	var req ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Geçersiz istek"})
	}
	if len(req.NewPassword) < 4 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Şifre en az 4 karakter olmalı"})
	}

	adminID := c.Get("user_id").(uint)
	if err := h.service.AdminResetPassword(adminID, uint(targetID), req.NewPassword); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Şifre sıfırlandı"})
}
