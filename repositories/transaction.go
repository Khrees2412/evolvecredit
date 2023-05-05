package repositories

import (
	"github.com/khrees2412/evolvecredit/database"
	"github.com/khrees2412/evolvecredit/models"
	"github.com/khrees2412/evolvecredit/types"
	"gorm.io/gorm"
	"time"
)

type ITransactionRepository interface {
	Create(transaction *models.Transaction) error
	FindById(userId string, transactionId string) (*models.Transaction, error)
	FindAllByUserId(userId string, entry types.TransactionEntry, status types.TransactionStatus, pagination types.Pagination) (*[]models.Transaction, error)
	FindAllByAccountId(accountId string) (*[]models.Transaction, error)
	FindAllByUserInRange(userId string, from time.Time) (*[]models.Transaction, error)
	Update(transaction *models.Transaction) error
	WithTx(tx *gorm.DB) ITransactionRepository
}

type transactionRepo struct {
	db *gorm.DB
}

// NewTransactionRepo will instantiate User Repository
func NewTransactionRepo() ITransactionRepository {
	return &transactionRepo{
		db: database.DB(),
	}
}

func (t *transactionRepo) Create(transaction *models.Transaction) error {
	return t.db.Create(transaction).Error
}

func (t *transactionRepo) FindById(userId string, transactionId string) (*models.Transaction, error) {

	var transaction models.Transaction
	if err := t.db.Where("user_id AND id = ?", userId, transactionId).First(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (t *transactionRepo) FindAllByUserId(userId string, entry types.TransactionEntry, status types.TransactionStatus, pagination types.Pagination) (*[]models.Transaction, error) {
	var transactions []models.Transaction

	chain := t.db.Scopes(paginate(pagination.Page, pagination.PageSize)).Where("user_id = ?", userId)

	if status != "" {
		chain = chain.Where("status = ?", status)
	}
	if entry != "" {
		chain = chain.Where("entry = ?", entry)
	}

	if err := chain.Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}

	return &transactions, nil
}

func (t *transactionRepo) FindAllByAccountId(accountId string) (*[]models.Transaction, error) {

	var transactions []models.Transaction
	if err := t.db.Where("account_id = ? ", accountId).Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}

	return &transactions, nil
}
func (t *transactionRepo) FindAllByUserInRange(userId string, from time.Time) (*[]models.Transaction, error) {

	var transactions []models.Transaction
	if err := t.db.Where("user_id = ? AND created_at >= ?", userId, from).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return &transactions, nil
}

func (t *transactionRepo) Update(transaction *models.Transaction) error {
	return t.db.Save(transaction).Error
}

func (t *transactionRepo) WithTx(tx *gorm.DB) ITransactionRepository {
	return &transactionRepo{db: tx}
}
