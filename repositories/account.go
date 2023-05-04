package repositories

import (
	"github.com/khrees2412/evolvecredit/database"
	"github.com/khrees2412/evolvecredit/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IAccountRepository interface {
	Create(account *models.Account) error
	Update(account *models.Account) error
	FindByAccountNumber(accountNumber string) (*models.Account, error)
	FindByUserId(userId string) (*models.Account, error)
	WithLock(locking *clause.Locking) IAccountRepository
	WithTx(tx *gorm.DB) IAccountRepository
	Delete(accountId string) error
}

type accountRepo struct {
	db *gorm.DB
}

// NewAccountRepo will instantiate Account Repository
func NewAccountRepo() IAccountRepository {
	return &accountRepo{
		db: database.DB(),
	}
}

func (r *accountRepo) Create(account *models.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepo) Update(account *models.Account) error {
	return r.db.Save(account).Error
}

func (r *accountRepo) FindByAccountNumber(accountNumber string) (*models.Account, error) {
	var account models.Account
	if err := r.db.Where("account_number = ?", accountNumber).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepo) FindByUserId(userId string) (*models.Account, error) {
	var account models.Account
	if err := r.db.Where("user_id = ?", userId).First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}
func (r *accountRepo) Delete(accountId string) error {
	var account models.Account
	if err := r.db.Where("id =  ?", accountId).Delete(&account).Error; err != nil {
		return err
	}
	return nil
}

func (r *accountRepo) WithLock(locking *clause.Locking) IAccountRepository {
	db := r.db.Clauses(locking)
	return &accountRepo{
		db: db,
	}
}

func (r *accountRepo) WithTx(tx *gorm.DB) IAccountRepository {
	return &accountRepo{db: tx}
}
