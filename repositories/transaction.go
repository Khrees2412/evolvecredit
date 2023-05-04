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
	FindAllByUserID(userId string, page int, pageSize int, status types.TransactionStatus, entry types.TransactionEntry) *[]models.Transaction
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

func (t *transactionRepo) FindByID(id uint) (*models.Transaction, error) {

	var transaction models.Transaction
	if err := t.db.Where("id = ?", id).Preload("Counterparty").First(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (t *transactionRepo) FindByReference(reference string) (*models.Transaction, error) {

	var transaction models.Transaction
	if err := t.db.Where("reference = ?", reference).First(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (t *transactionRepo) FindByExRef(reference string) (*models.Transaction, error) {

	var transaction models.Transaction
	if err := t.db.Where("external_ref = ?", reference).First(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (t *transactionRepo) FindAllByUserID(userId string, page int, pageSize int, status types.TransactionStatus, entry types.TransactionEntry) *[]models.Transaction {

	var transactions []models.Transaction

	chain := t.db.Scopes(paginate(page, pageSize)).Where("user_id = ?", userId)

	if status != "" {
		chain = chain.Where("status = ?", status)
	}
	if entry != "" {
		chain = chain.Where("entry = ?", entry)
	}

	chain.Order("id DESC").Find(&transactions)

	return &transactions
}

func (t *transactionRepo) FindAllByAccountId(accountId string) (*[]models.Transaction, error) {

	var transactions []models.Transaction
	if err := t.db.Where("account_id = ? ", accountId).Order("id DESC").Find(&transactions).Error; err != nil {
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
