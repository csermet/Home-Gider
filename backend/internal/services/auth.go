package services

import (
	"errors"
	"time"

	"github.com/caner/home-gider/internal/middleware"
	"github.com/caner/home-gider/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db     *gorm.DB
	secret string
}

func NewAuthService(db *gorm.DB, secret string) *AuthService {
	return &AuthService{db: db, secret: secret}
}

func (s *AuthService) Login(username, password string) (*models.User, string, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", err
	}

	claims := &middleware.JWTClaims{
		UserID:      user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		IsAdmin:     user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return nil, "", err
	}

	return &user, tokenStr, nil
}

func (s *AuthService) GetUser(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) ChangePassword(userID uint, newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"password_hash":        string(hash),
		"must_change_password": false,
	}).Error
}

// AdminResetPassword admin kullanıcının başka bir kullanıcının şifresini sıfırlaması
func (s *AuthService) AdminResetPassword(adminID, targetUserID uint, newPassword string) error {
	var admin models.User
	if err := s.db.First(&admin, adminID).Error; err != nil {
		return err
	}
	if !admin.IsAdmin {
		return errors.New("bu işlem için admin yetkisi gerekiyor")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.db.Model(&models.User{}).Where("id = ?", targetUserID).Updates(map[string]interface{}{
		"password_hash":        string(hash),
		"must_change_password": true,
	}).Error
}

// ListUsers admin için tüm kullanıcıları listele
func (s *AuthService) ListUsers() ([]models.User, error) {
	var users []models.User
	err := s.db.Order("id ASC").Find(&users).Error
	return users, err
}
