package database

import (
	"log"

	"github.com/caner/home-gider/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	seedUsers(db)
	seedCategories(db)
}

func seedUsers(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	// Admin kullanıcı — ilk girişte şifre belirleyecek
	adminHash, _ := bcrypt.GenerateFromPassword([]byte("temp"), bcrypt.DefaultCost)
	db.Create(&models.User{
		Username:           "admin",
		PasswordHash:       string(adminHash),
		DisplayName:        "Admin",
		IsAdmin:            true,
		MustChangePassword: true,
	})

	// CNR ve CNS — ilk girişte şifre belirleyecekler
	// Geçici şifre: "temp" (sadece ilk giriş için)
	tempHash, _ := bcrypt.GenerateFromPassword([]byte("temp"), bcrypt.DefaultCost)
	db.Create(&models.User{
		Username:           "cnr",
		PasswordHash:       string(tempHash),
		DisplayName:        "CNR",
		IsAdmin:            false,
		MustChangePassword: true,
	})
	db.Create(&models.User{
		Username:           "cns",
		PasswordHash:       string(tempHash),
		DisplayName:        "CNS",
		IsAdmin:            false,
		MustChangePassword: true,
	})

	log.Println("Kullanıcılar oluşturuldu (CNR, CNS, Admin)")
}

func seedCategories(db *gorm.DB) {
	var count int64
	db.Model(&models.Category{}).Count(&count)
	if count > 0 {
		return
	}

	categories := []models.Category{
		{Name: "Kira", Icon: "building"},
		{Name: "Fatura", Icon: "receipt"},
		{Name: "Market", Icon: "shopping-cart"},
		{Name: "Ev", Icon: "home"},
		{Name: "Eğlence", Icon: "gamepad"},
		{Name: "Ulaşım", Icon: "car"},
		{Name: "Sağlık", Icon: "heart-pulse"},
		{Name: "Giyim", Icon: "shirt"},
		{Name: "Yemek", Icon: "utensils"},
		{Name: "Diğer", Icon: "ellipsis"},
	}

	db.Create(&categories)
	log.Println("Kategoriler oluşturuldu")
}
