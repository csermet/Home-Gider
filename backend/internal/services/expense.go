package services

import (
	"errors"
	"time"

	"github.com/caner/home-gider/internal/models"
	"gorm.io/gorm"
)

type ExpenseService struct {
	db *gorm.DB
}

func NewExpenseService(db *gorm.DB) *ExpenseService {
	return &ExpenseService{db: db}
}

func (s *ExpenseService) List(month, year int) ([]models.Expense, error) {
	var expenses []models.Expense
	err := s.db.
		Preload("Creator").
		Preload("Category").
		Preload("Approver").
		Where("expense_month = ? AND expense_year = ?", month, year).
		Order("expense_date DESC, created_at DESC").
		Find(&expenses).Error
	return expenses, err
}

func (s *ExpenseService) Create(expense *models.Expense) error {
	expense.ExpenseMonth = int(expense.ExpenseDate.Month())
	expense.ExpenseYear = expense.ExpenseDate.Year()

	if !expense.IsShared {
		expense.Status = models.StatusApproved
	}

	return s.db.Create(expense).Error
}

func (s *ExpenseService) Update(id, userID uint, updates map[string]interface{}) error {
	var expense models.Expense
	if err := s.db.First(&expense, id).Error; err != nil {
		return err
	}
	if expense.CreatedBy != userID {
		return errors.New("sadece kendi eklediğiniz giderleri düzenleyebilirsiniz")
	}
	if expense.Status != models.StatusPending {
		return errors.New("sadece onay bekleyen giderler düzenlenebilir")
	}
	return s.db.Model(&expense).Updates(updates).Error
}

func (s *ExpenseService) Delete(id, userID uint) error {
	var expense models.Expense
	if err := s.db.First(&expense, id).Error; err != nil {
		return err
	}
	if expense.CreatedBy != userID {
		return errors.New("sadece kendi eklediğiniz giderleri silebilirsiniz")
	}
	if expense.Status != models.StatusPending {
		return errors.New("sadece onay bekleyen giderler silinebilir")
	}
	return s.db.Delete(&expense).Error
}

func (s *ExpenseService) Approve(id, approverID uint, isAdmin bool) error {
	var expense models.Expense
	if err := s.db.First(&expense, id).Error; err != nil {
		return err
	}
	if !isAdmin && expense.CreatedBy == approverID {
		return errors.New("kendi eklediğiniz gideri onaylayamazsınız")
	}
	if expense.Status != models.StatusPending {
		return errors.New("bu gider zaten işlenmiş")
	}
	now := time.Now()
	return s.db.Model(&expense).Updates(map[string]interface{}{
		"status":      models.StatusApproved,
		"approved_by": approverID,
		"approved_at": now,
	}).Error
}

func (s *ExpenseService) Reject(id, userID uint, isAdmin bool) error {
	var expense models.Expense
	if err := s.db.First(&expense, id).Error; err != nil {
		return err
	}
	if !isAdmin && expense.CreatedBy == userID {
		return errors.New("kendi eklediğiniz gideri reddedemezsiniz")
	}
	if expense.Status != models.StatusPending {
		return errors.New("bu gider zaten işlenmiş")
	}
	return s.db.Model(&expense).Update("status", models.StatusRejected).Error
}

func (s *ExpenseService) GetByID(id uint) (*models.Expense, error) {
	var expense models.Expense
	err := s.db.Preload("Creator").Preload("Category").Preload("Approver").First(&expense, id).Error
	return &expense, err
}
