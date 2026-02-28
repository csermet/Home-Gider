package models

import "time"

type RecurringType string

const (
	TypeInstallment RecurringType = "installment"
	TypeRecurring   RecurringType = "recurring"
)

type RecurringExpense struct {
	ID                    uint          `json:"id" gorm:"primaryKey"`
	CreatedBy             uint          `json:"created_by" gorm:"not null"`
	Creator               User          `json:"creator" gorm:"foreignKey:CreatedBy"`
	CategoryID            uint          `json:"category_id" gorm:"not null"`
	Category              Category      `json:"category" gorm:"foreignKey:CategoryID"`
	Description           string        `json:"description" gorm:"size:255;not null"`
	Amount                float64       `json:"amount" gorm:"type:decimal(10,2);not null"`
	TotalAmount           *float64      `json:"total_amount" gorm:"type:decimal(10,2)"`
	Type                  RecurringType `json:"type" gorm:"size:20;not null"`
	InstallmentCount      *int          `json:"installment_count"`
	InstallmentsRemaining *int          `json:"installments_remaining"`
	IsShared              bool          `json:"is_shared" gorm:"default:true"`
	SplitRatio            float64       `json:"split_ratio" gorm:"type:decimal(5,2);default:50"`
	IsActive              bool          `json:"is_active" gorm:"default:true"`
	Status                ExpenseStatus `json:"status" gorm:"size:20;default:'pending'"`
	ApprovedBy            *uint         `json:"approved_by"`
	Approver              *User         `json:"approver,omitempty" gorm:"foreignKey:ApprovedBy"`
	CreatedAt             time.Time     `json:"created_at"`
}
