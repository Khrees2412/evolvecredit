package services

import (
	"errors"
	"github.com/khrees2412/evolvecredit/models"
	"github.com/khrees2412/evolvecredit/repositories"
	"github.com/khrees2412/evolvecredit/types"
)

type IPaymentService interface {
	Deposit(user models.User, request *types.DepositRequest) error
	Withdrawal(user *models.User, request *types.WithdrawalRequest) error
}

type paymentService struct {
	accountRepo     repositories.IAccountRepository
	transactionRepo repositories.ITransactionRepository
	userRepo        repositories.IUserRepository
	ledgerRepo      repositories.IWalletLedgerRepository
	accountService  IAccountService
}

func NewPaymentService() IPaymentService {
	return &paymentService{
		accountRepo:     repositories.NewAccountRepo(),
		transactionRepo: repositories.NewTransactionRepo(),
		ledgerRepo:      repositories.NewWalletLedgerRepo(),
		userRepo:        repositories.NewUserRepo(),
		accountService:  NewAccountService(),
	}
}

var (
	ErrAccountIsDisabled = errors.New("this user's account is currently disabled")
)

func (p paymentService) Deposit(user models.User, request *types.DepositRequest) error {
	uw := repositories.NewGormUnitOfWork()
	tx, err := uw.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	account, err := p.accountRepo.FindByUserId(user.UserId)
	if !account.IsActive {
		return ErrAccountIsDisabled
	}
	if err != nil {
		return err
	}

	response, err := p.accountService.CreditAccount(user.UserId, request.Amount, tx)
	if err != nil {
		return err
	}

	trx := &models.Transaction{
		UserId:        user.UserId,
		AccountNumber: request.AccountNumber,
		Status:        types.Success,
		Entry:         types.Credit,
		Type:          types.Deposit,
		Amount:        request.Amount,
		Balance:       response.CurrentBalance,
		Summary:       request.Reason,
	}

	err = p.transactionRepo.WithTx(tx).Create(trx)

	if err != nil {
		return err
	}

	err = p.ledgerRepo.WithTx(tx).Create(&models.WalletLedger{
		TransactionId:   trx.Id,
		AccountNumber:   response.AccountNumber,
		Entry:           types.Credit,
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

func (p paymentService) Withdrawal(user *models.User, request *types.WithdrawalRequest) error {
	uw := repositories.NewGormUnitOfWork()
	tx, err := uw.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	account, err := p.accountRepo.FindByUserId(user.UserId)
	if !account.IsActive {
		return ErrAccountIsDisabled
	}
	response, err := p.accountService.DebitAccount(user.UserId, request.Amount, tx)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	trx := &models.Transaction{
		UserId:  user.UserId,
		Status:  types.Success,
		Entry:   types.Debit,
		Type:    types.Withdrawal,
		Amount:  request.Amount,
		Summary: request.Reason,
		Balance: response.CurrentBalance,
	}

	err = p.transactionRepo.WithTx(tx).Create(trx)

	if err != nil {
		return err
	}

	err = uw.Commit(tx)

	if err != nil {
		return err
	}

	return err
}
