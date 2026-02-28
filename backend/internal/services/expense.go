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
		Preload("DeleteRequester").
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

func (s *ExpenseService) Delete(id, userID uint, isAdmin bool) error {
	var expense models.Expense
	if err := s.db.First(&expense, id).Error; err != nil {
		return err
	}
	// Admin direkt silebilir
	if isAdmin {
		return s.db.Delete(&expense).Error
	}
	// Normal kullanıcı: silme talep et
	if expense.DeleteRequestedBy != nil {
		return errors.New("bu gider için zaten silme talebi var")
	}
	return s.db.Model(&expense).Update("delete_requested_by", userID).Error
}

func (s *ExpenseService) ConfirmDelete(id, userID uint) error {
	var expense models.Expense
	if err := s.db.First(&expense, id).Error; err != nil {
		return err
	}
	if expense.DeleteRequestedBy == nil {
		return errors.New("bu gider için silme talebi yok")
	}
	if *expense.DeleteRequestedBy == userID {
		return errors.New("kendi silme talebinizi onaylayamazsınız")
	}
	return s.db.Delete(&expense).Error
}

func (s *ExpenseService) CancelDelete(id, userID uint) error {
	var expense models.Expense
	if err := s.db.First(&expense, id).Error; err != nil {
		return err
	}
	if expense.DeleteRequestedBy == nil {
		return errors.New("bu gider için silme talebi yok")
	}
	if *expense.DeleteRequestedBy != userID {
		return errors.New("sadece talebi oluşturan kişi iptal edebilir")
	}
	return s.db.Model(&expense).Update("delete_requested_by", nil).Error
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
	err := s.db.Preload("Creator").Preload("Category").Preload("Approver").Preload("DeleteRequester").First(&expense, id).Error
	return &expense, err
}
