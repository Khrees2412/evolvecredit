package services

import (
	"github.com/khrees2412/evolvecredit/models"
	"github.com/khrees2412/evolvecredit/repositories"
	"github.com/khrees2412/evolvecredit/types"
)

type ITransactionService interface {
	GetTransaction(userId string, transactionId string) (*models.Transaction, error)
	GetAllTransactions(userId string, entry types.TransactionEntry, status types.TransactionStatus, pagination types.Pagination) (*[]models.Transaction, error)
}

type transactionService struct {
	transactionRepo repositories.ITransactionRepository
}

func NewTransactionService() ITransactionService {
	return &transactionService{transactionRepo: repositories.NewTransactionRepo()}
}

func (ts transactionService) GetAllTransactions(userId string, entry types.TransactionEntry, status types.TransactionStatus, pagination types.Pagination) (*[]models.Transaction, error) {
	transactions, err := ts.transactionRepo.FindAllByUserId(userId, entry, status, pagination)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (ts transactionService) GetTransaction(userId string, transactionId string) (*models.Transaction, error) {
	transaction, err := ts.transactionRepo.FindById(userId, transactionId)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}
