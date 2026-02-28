package models

import "time"

// Payment kısmi ödeme kaydı — ay içinde birden fazla ödeme yapılabilir
type Payment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Month     int       `json:"month" gorm:"not null"`
	Year      int       `json:"year" gorm:"not null"`
	PayerID   uint      `json:"payer_id" gorm:"not null"`
	Payer     User      `json:"payer" gorm:"foreignKey:PayerID"`
	PayeeID   uint      `json:"payee_id" gorm:"not null"`
	Payee     User      `json:"payee" gorm:"foreignKey:PayeeID"`
	Amount    float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	CreatedAt time.Time `json:"created_at"`
}
