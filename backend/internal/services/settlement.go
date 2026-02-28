package services

import (
	"errors"
	"math"

	"github.com/caner/home-gider/internal/models"
	"gorm.io/gorm"
)

type SettlementService struct {
	db *gorm.DB
}

func NewSettlementService(db *gorm.DB) *SettlementService {
	return &SettlementService{db: db}
}

type UserSummary struct {
	UserID      uint    `json:"user_id"`
	DisplayName string  `json:"display_name"`
	TotalPaid   float64 `json:"total_paid"`
	TotalShare  float64 `json:"total_share"`
	Balance     float64 `json:"balance"`
}

type MonthlySummary struct {
	Month             int           `json:"month"`
	Year              int           `json:"year"`
	TotalExpenses     float64       `json:"total_expenses"`
	SharedExpenses    float64       `json:"shared_expenses"`
	UserSummaries     []UserSummary `json:"user_summaries"`
	DebtorID          *uint         `json:"debtor_id"`
	CreditorID        *uint         `json:"creditor_id"`
	DebtAmount        float64       `json:"debt_amount"`
	TotalPayments     float64       `json:"total_payments"`
	RemainingDebt     float64       `json:"remaining_debt"`
	CategoryBreakdown []CategorySum `json:"category_breakdown"`
}

type CategorySum struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	CategoryIcon string  `json:"category_icon"`
	Total        float64 `json:"total"`
}

func (s *SettlementService) GetMonthlySummary(month, year int, sharedOnly bool) (*MonthlySummary, error) {
	var expenses []models.Expense
	query := s.db.Preload("Creator").Preload("Category").
		Where("expense_month = ? AND expense_year = ? AND status = ?",
			month, year, models.StatusApproved)
	if sharedOnly {
		query = query.Where("is_shared = ?", true)
	}
	err := query.Find(&expenses).Error
	if err != nil {
		return nil, err
	}

	// Sadece admin olmayan kullanıcıları al (hesaplamaya admin dahil değil)
	var users []models.User
	s.db.Where("is_admin = ?", false).Find(&users)

	userMap := make(map[uint]*UserSummary)
	for _, u := range users {
		userMap[u.ID] = &UserSummary{
			UserID:      u.ID,
			DisplayName: u.DisplayName,
		}
	}

	var totalExpenses, sharedExpenses float64
	categoryMap := make(map[uint]*CategorySum)

	for _, e := range expenses {
		totalExpenses += e.Amount

		if _, ok := categoryMap[e.CategoryID]; !ok {
			categoryMap[e.CategoryID] = &CategorySum{
				CategoryID:   e.CategoryID,
				CategoryName: e.Category.Name,
				CategoryIcon: e.Category.Icon,
			}
		}
		categoryMap[e.CategoryID].Total += e.Amount

		if summary, ok := userMap[e.CreatedBy]; ok {
			summary.TotalPaid += e.Amount
		}

		if !e.IsShared {
			if summary, ok := userMap[e.CreatedBy]; ok {
				summary.TotalShare += e.Amount
			}
			continue
		}

		sharedExpenses += e.Amount

		creatorShare := e.Amount * e.SplitRatio / 100
		otherShare := e.Amount - creatorShare

		if summary, ok := userMap[e.CreatedBy]; ok {
			summary.TotalShare += creatorShare
		}

		for uid, summary := range userMap {
			if uid != e.CreatedBy {
				summary.TotalShare += otherShare
			}
		}
	}

	summaries := make([]UserSummary, 0, len(userMap))
	for _, summary := range userMap {
		summary.Balance = math.Round((summary.TotalPaid-summary.TotalShare)*100) / 100
		summaries = append(summaries, *summary)
	}

	result := &MonthlySummary{
		Month:          month,
		Year:           year,
		TotalExpenses:  math.Round(totalExpenses*100) / 100,
		SharedExpenses: math.Round(sharedExpenses*100) / 100,
		UserSummaries:  summaries,
	}

	// Kim kime borçlu?
	if len(summaries) == 2 {
		if summaries[0].Balance > 0 {
			result.CreditorID = &summaries[0].UserID
			result.DebtorID = &summaries[1].UserID
			result.DebtAmount = math.Round(summaries[0].Balance*100) / 100
		} else if summaries[1].Balance > 0 {
			result.CreditorID = &summaries[1].UserID
			result.DebtorID = &summaries[0].UserID
			result.DebtAmount = math.Round(summaries[1].Balance*100) / 100
		}
	}

	// Yapılan ödemeleri hesapla
	var payments []models.Payment
	s.db.Where("month = ? AND year = ?", month, year).Find(&payments)
	var totalPayments float64
	for _, p := range payments {
		totalPayments += p.Amount
	}
	result.TotalPayments = math.Round(totalPayments*100) / 100
	result.RemainingDebt = math.Round((result.DebtAmount-totalPayments)*100) / 100
	if result.RemainingDebt < 0 {
		result.RemainingDebt = 0
	}

	categories := make([]CategorySum, 0, len(categoryMap))
	for _, cs := range categoryMap {
		cs.Total = math.Round(cs.Total*100) / 100
		categories = append(categories, *cs)
	}
	result.CategoryBreakdown = categories

	return result, nil
}

func (s *SettlementService) GetPayments(month, year int) ([]models.Payment, error) {
	var payments []models.Payment
	err := s.db.Preload("Payer").Preload("Payee").
		Where("month = ? AND year = ?", month, year).
		Order("created_at DESC").
		Find(&payments).Error
	return payments, err
}

func (s *SettlementService) AddPayment(month, year int, payerID, payeeID uint, amount float64) error {
	// Borç miktarını hesapla
	summary, err := s.GetMonthlySummary(month, year, true)
	if err != nil {
		return err
	}
	if summary.RemainingDebt <= 0 {
		return errors.New("bu ay için borç kalmamış")
	}
	if amount > summary.RemainingDebt {
		return errors.New("ödeme tutarı kalan borçtan büyük olamaz")
	}

	payment := models.Payment{
		Month:   month,
		Year:    year,
		PayerID: payerID,
		PayeeID: payeeID,
		Amount:  amount,
	}
	return s.db.Create(&payment).Error
}

func (s *SettlementService) DeletePayment(id uint) error {
	return s.db.Delete(&models.Payment{}, id).Error
}
