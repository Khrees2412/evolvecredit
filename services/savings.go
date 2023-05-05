package services

import (
	"github.com/khrees2412/evolvecredit/models"
	"github.com/khrees2412/evolvecredit/repositories"
	"github.com/khrees2412/evolvecredit/types"
)

type ISavingsService interface {
	GetSavings(accountNumber string) (*types.SavingsResponse, error)
	SaveFunds(userId string, request types.SavingsRequest) error
}

type savingsService struct {
	accountRepo     repositories.IAccountRepository
	savingsRepo     repositories.ISavingsRepository
	transactionRepo repositories.ITransactionRepository
	ledgerRepo      repositories.IWalletLedgerRepository
	accountService  IAccountService
}

func NewSavingsService() ISavingsService {
	return &savingsService{
		accountRepo:     repositories.NewAccountRepo(),
		savingsRepo:     repositories.NewSavingsRepo(),
		transactionRepo: repositories.NewTransactionRepo(),
		ledgerRepo:      repositories.NewWalletLedgerRepo(),
		accountService:  NewAccountService(),
	}
}

func (s savingsService) GetSavings(accountNumber string) (*types.SavingsResponse, error) {
	savings, err := s.savingsRepo.FindByAccountNumber(accountNumber)
	if err != nil {
		return nil, err
	}
	account, err := s.accountRepo.FindByAccountNumber(accountNumber)
	return &types.SavingsResponse{
		AccountNumber:  accountNumber,
		CurrentBalance: account.AvailableBalance,
		LockedAmount:   savings.Amount,
	}, nil
}

func (s savingsService) SaveFunds(userId string, request types.SavingsRequest) error {

	uw := repositories.NewGormUnitOfWork()
	tx, err := uw.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	account, err := s.accountRepo.FindByAccountNumber(request.AccountNumber)
	if !account.IsActive {
		return ErrAccountIsDisabled
	}
	if err != nil {
		return err
	}
	response, err := s.accountService.LockFunds(userId, request.AccountNumber, request.Amount, request.Duration, tx)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	trx := &models.Transaction{
		UserId:        userId,
		AccountNumber: request.AccountNumber,
		Status:        types.Success,
		Amount:        request.Amount,
		Balance:       response.CurrentBalance,
		LockedAmount:  response.LockedAmount,
		Summary:       "Savings",
	}

	err = s.transactionRepo.WithTx(tx).Create(trx)

	if err != nil {
		return err
	}

	err = s.ledgerRepo.WithTx(tx).Create(&models.WalletLedger{
		TransactionId:   trx.Id,
		AccountNumber:   response.AccountNumber,
		PreviousBalance: response.PreviousBalance,
		CurrentBalance:  response.CurrentBalance,
	})

	if err != nil {
		return err
	}

	err = uw.Commit(tx)

	if err != nil {
		return err
	}
	return nil
}
