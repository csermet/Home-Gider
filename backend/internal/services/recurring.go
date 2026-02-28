package services

import (
	"errors"
	"time"

	"github.com/caner/home-gider/internal/models"
	"gorm.io/gorm"
)

type RecurringService struct {
	db *gorm.DB
}

func NewRecurringService(db *gorm.DB) *RecurringService {
	return &RecurringService{db: db}
}

func (s *RecurringService) List() ([]models.RecurringExpense, error) {
	var items []models.RecurringExpense
	err := s.db.Preload("Creator").Preload("Category").Preload("Approver").
		Order("created_at DESC").Find(&items).Error
	return items, err
}

func (s *RecurringService) Create(item *models.RecurringExpense) error {
	if item.Type == models.TypeInstallment {
		if item.InstallmentCount == nil || *item.InstallmentCount <= 0 {
			return errors.New("taksit sayısı belirtilmelidir")
		}
		remaining := *item.InstallmentCount
		item.InstallmentsRemaining = &remaining
	}
	return s.db.Create(item).Error
}

func (s *RecurringService) Update(id, userID uint, updates map[string]interface{}) error {
	var item models.RecurringExpense
	if err := s.db.First(&item, id).Error; err != nil {
		return err
	}
	if item.CreatedBy != userID {
		return errors.New("sadece kendi oluşturduğunuz şablonları düzenleyebilirsiniz")
	}
	return s.db.Model(&item).Updates(updates).Error
}

func (s *RecurringService) Delete(id, userID uint) error {
	var item models.RecurringExpense
	if err := s.db.First(&item, id).Error; err != nil {
		return err
	}
	if item.CreatedBy != userID {
		return errors.New("sadece kendi oluşturduğunuz şablonları silebilirsiniz")
	}
	return s.db.Model(&item).Update("is_active", false).Error
}

func (s *RecurringService) Approve(id, approverID uint, isAdmin bool) error {
	var item models.RecurringExpense
	if err := s.db.First(&item, id).Error; err != nil {
		return err
	}
	if !isAdmin && item.CreatedBy == approverID {
		return errors.New("kendi oluşturduğunuz şablonu onaylayamazsınız")
	}
	if item.Status != models.StatusPending {
		return errors.New("bu şablon zaten işlenmiş")
	}

	if err := s.db.Model(&item).Updates(map[string]interface{}{
		"status":      models.StatusApproved,
		"approved_by": approverID,
	}).Error; err != nil {
		return err
	}

	// Onaylandığında hemen bu ay için gider oluştur
	return s.createExpenseForMonth(&item, time.Now())
}

func (s *RecurringService) Reject(id, approverID uint, isAdmin bool) error {
	var item models.RecurringExpense
	if err := s.db.First(&item, id).Error; err != nil {
		return err
	}
	if !isAdmin && item.CreatedBy == approverID {
		return errors.New("kendi oluşturduğunuz şablonu reddedemezsiniz")
	}
	if item.Status != models.StatusPending {
		return errors.New("bu şablon zaten işlenmiş")
	}
	return s.db.Model(&item).Update("status", models.StatusRejected).Error
}

// createExpenseForMonth verilen ay için gider kaydı oluşturur
func (s *RecurringService) createExpenseForMonth(item *models.RecurringExpense, date time.Time) error {
	month := int(date.Month())
	year := date.Year()

	// Bu ay için zaten kayıt var mı?
	var count int64
	s.db.Model(&models.Expense{}).
		Where("recurring_expense_id = ? AND expense_month = ? AND expense_year = ?",
			item.ID, month, year).
		Count(&count)
	if count > 0 {
		return nil
	}

	// Taksit kontrolü
	if item.Type == models.TypeInstallment {
		if item.InstallmentsRemaining != nil && *item.InstallmentsRemaining <= 0 {
			s.db.Model(item).Update("is_active", false)
			return nil
		}
	}

	installmentNo := 0
	var installmentTotal *int
	if item.Type == models.TypeInstallment && item.InstallmentCount != nil && item.InstallmentsRemaining != nil {
		installmentNo = *item.InstallmentCount - *item.InstallmentsRemaining + 1
		installmentTotal = item.InstallmentCount
	}

	expense := models.Expense{
		CreatedBy:          item.CreatedBy,
		CategoryID:         item.CategoryID,
		Description:        item.Description,
		Amount:             item.Amount,
		ExpenseDate:        date,
		ExpenseMonth:       month,
		ExpenseYear:        year,
		IsShared:           item.IsShared,
		SplitRatio:         item.SplitRatio,
		IsInstallment:      item.Type == models.TypeInstallment,
		RecurringExpenseID: &item.ID,
		Status:             models.StatusApproved,
	}
	if installmentNo > 0 {
		expense.InstallmentNo = &installmentNo
		expense.InstallmentTotal = installmentTotal
	}

	if err := s.db.Create(&expense).Error; err != nil {
		return err
	}

	// Taksitte kalan sayıyı azalt
	if item.Type == models.TypeInstallment && item.InstallmentsRemaining != nil {
		newRemaining := *item.InstallmentsRemaining - 1
		updates := map[string]interface{}{"installments_remaining": newRemaining}
		if newRemaining <= 0 {
			updates["is_active"] = false
		}
		s.db.Model(item).Updates(updates)
	}

	return nil
}

// ProcessRecurring her gün çalışır — bu ay için eksik gider kayıtlarını oluşturur
func (s *RecurringService) ProcessRecurring() error {
	now := time.Now()

	var items []models.RecurringExpense
	err := s.db.Where("is_active = ? AND status = ?",
		true, models.StatusApproved).Find(&items).Error
	if err != nil {
		return err
	}

	for _, item := range items {
		if err := s.createExpenseForMonth(&item, now); err != nil {
			return err
		}
	}

	return nil
}
