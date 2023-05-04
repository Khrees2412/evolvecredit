package repositories

import (
	"github.com/khrees2412/evolvecredit/database"
	"github.com/khrees2412/evolvecredit/models"
	"gorm.io/gorm"
)

type IAccountRepository interface {
	Create(account *models.Account) error
	Update(account *models.Account) error
	FindByAccountId(accountId string) (*models.Account, error)
	FindAllByUserId(userId string) (*[]models.Account, error)
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

func (r *accountRepo) FindByAccountId(accountId string) (*models.Account, error) {
	var account models.Account
	if err := r.db.Where("id = ?", accountId).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *accountRepo) FindAllByUserId(userId string) (*[]models.Account, error) {
	var account []models.Account
	if err := r.db.Where("user_id = ?", userId).Find(&account).Error; err != nil {
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
