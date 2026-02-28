package models

import "time"

type ExpenseStatus string

const (
	StatusPending  ExpenseStatus = "pending"
	StatusApproved ExpenseStatus = "approved"
	StatusRejected ExpenseStatus = "rejected"
)

type Expense struct {
	ID                 uint          `json:"id" gorm:"primaryKey"`
	CreatedBy          uint          `json:"created_by" gorm:"not null"`
	Creator            User          `json:"creator" gorm:"foreignKey:CreatedBy"`
	CategoryID         uint          `json:"category_id" gorm:"not null"`
	Category           Category      `json:"category" gorm:"foreignKey:CategoryID"`
	Description        string        `json:"description" gorm:"size:255;not null"`
	Amount             float64       `json:"amount" gorm:"type:decimal(10,2);not null"`
	ExpenseDate        time.Time     `json:"expense_date" gorm:"type:date;not null"`
	ExpenseMonth       int           `json:"expense_month" gorm:"not null"`
	ExpenseYear        int           `json:"expense_year" gorm:"not null"`
	IsShared           bool          `json:"is_shared" gorm:"default:true"`
	SplitRatio         float64       `json:"split_ratio" gorm:"type:decimal(5,2);default:50"`
	IsInstallment      bool          `json:"is_installment" gorm:"default:false"`
	InstallmentNo      *int          `json:"installment_no"`
	InstallmentTotal   *int          `json:"installment_total"`
	RecurringExpenseID *uint         `json:"recurring_expense_id"`
	Status             ExpenseStatus `json:"status" gorm:"size:20;default:'pending'"`
	ApprovedBy         *uint         `json:"approved_by"`
	Approver           *User         `json:"approver,omitempty" gorm:"foreignKey:ApprovedBy"`
	ApprovedAt         *time.Time    `json:"approved_at"`
	DeleteRequestedBy  *uint         `json:"delete_requested_by"`
	DeleteRequester    *User         `json:"delete_requester,omitempty" gorm:"foreignKey:DeleteRequestedBy"`
	CreatedAt          time.Time     `json:"created_at"`
}
