package repositories

import (
	"github.com/khrees2412/evolvecredit/database"
	"github.com/khrees2412/evolvecredit/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ISavingsRepository interface {
	Create(savings *models.Savings) error
	Update(savings *models.Savings) error
	FindByAccountNumber(accountNumber string) (*models.Savings, error)
	FindByUserId(userId string) (*models.Savings, error)
	WithLock(locking *clause.Locking) ISavingsRepository
	WithTx(tx *gorm.DB) ISavingsRepository
	Delete(savingsId string) error
}

type savingsRepo struct {
	db *gorm.DB
}

// NewSavingsRepo will instantiate Savings Repository
func NewSavingsRepo() ISavingsRepository {
	return &savingsRepo{
		db: database.DB(),
	}
}

func (r *savingsRepo) Create(savings *models.Savings) error {
	return r.db.Create(savings).Error
}

func (r *savingsRepo) Update(savings *models.Savings) error {
	return r.db.Save(savings).Error
}

func (r *savingsRepo) FindByAccountNumber(accountNumber string) (*models.Savings, error) {
	var savings models.Savings
	if err := r.db.Where("account_number = ?", accountNumber).First(&savings).Error; err != nil {
		return nil, err
	}
	return &savings, nil
}

func (r *savingsRepo) FindByUserId(userId string) (*models.Savings, error) {
	var savings models.Savings
	if err := r.db.Where("user_id = ?", userId).First(&savings).Error; err != nil {
		return nil, err
	}

	return &savings, nil
}
func (r *savingsRepo) Delete(accountNumber string) error {
	var savings models.Savings
	if err := r.db.Where("account_number = ?", accountNumber).Delete(&savings).Error; err != nil {
		return err
	}
	return nil
}

func (r *savingsRepo) WithLock(locking *clause.Locking) ISavingsRepository {
	db := r.db.Clauses(locking)
	return &savingsRepo{
		db: db,
	}
}

func (r *savingsRepo) WithTx(tx *gorm.DB) ISavingsRepository {
	return &savingsRepo{db: tx}
}
