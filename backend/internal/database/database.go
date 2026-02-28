package database

import (
	"log"

	"github.com/caner/home-gider/internal/config"
	"github.com/caner/home-gider/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Expense{},
		&models.RecurringExpense{},
		&models.Payment{},
	)
	if err != nil {
		log.Fatalf("Migration hatası: %v", err)
	}

	log.Println("Veritabanı bağlantısı başarılı")
	return db
}
