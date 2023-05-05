package services

import (
	"errors"
	"github.com/khrees2412/evolvecredit/models"
	"github.com/khrees2412/evolvecredit/repositories"
	"github.com/khrees2412/evolvecredit/types"
	"github.com/khrees2412/evolvecredit/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IAccountService interface {
	GetAccount(userId string) (*models.Account, error)
	CreateAccount(userId string) error
	CreditAccount(accountNumber string, amount int64, trx *gorm.DB) (*types.AccountResponse, error)
	DebitAccount(accountNumber string, amount int64, trx *gorm.DB) (*types.AccountResponse, error)
	LockFunds(userId string, accountNumber string, amount int64, duration int, trx *gorm.DB) (*types.SavingsResponse, error)
}

type accountService struct {
	userRepo    repositories.IUserRepository
	accountRepo repositories.IAccountRepository
	savingsRepo repositories.ISavingsRepository
}

func NewAccountService() IAccountService {
	return &accountService{
		userRepo:    repositories.NewUserRepo(),
		accountRepo: repositories.NewAccountRepo(),
		savingsRepo: repositories.NewSavingsRepo(),
	}
}

var (
	ErrAccountAlreadyExists = errors.New("an account already exists for this user")
)

func (as accountService) GetAccount(userId string) (*models.Account, error) {
	account, err := as.accountRepo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (as accountService) CreateAccount(userId string) error {
	_, err := as.accountRepo.FindByUserId(userId)
	if err != nil {
		return ErrAccountAlreadyExists
	}
	acct := models.Account{
		UserId:           userId,
		AccountNumber:    utils.GenerateAccountNumber(),
		LedgerBalance:    0,
		AvailableBalance: 0,
		TotalLocked:      0,
		IsActive:         true,
	}
	err = as.accountRepo.Create(&acct)
	if err != nil {
		return err
	}
	return nil
}

func (as accountService) CreditAccount(accountNumber string, amount int64, trx *gorm.DB) (*types.AccountResponse, error) {

	var lockedDB repositories.IAccountRepository

	if trx == nil {
		lockedDB = as.accountRepo.WithLock(&clause.Locking{
			Strength: "UPDATE",
		})
	} else {
		lockedDB = as.accountRepo.WithTx(trx).WithLock(&clause.Locking{
			Strength: "UPDATE",
		})
	}

	account, err := lockedDB.FindByAccountNumber(accountNumber)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			account = &models.Account{
				AccountNumber: accountNumber,
			}
			err = lockedDB.Create(account)
			if err != nil {
				return nil, errors.New("an error occurred when retrieving account, please try again")
			}
		} else {
			return nil, errors.New("an error occurred when retrieving account, please try again")
		}
	}

	prevBal := account.AvailableBalance
	account.AvailableBalance += amount
	err = lockedDB.Update(account)

	if err != nil {
		return nil, err
	}

	res := &types.AccountResponse{
		CurrentBalance:  account.AvailableBalance,
		PreviousBalance: prevBal,
		AccountNumber:   account.AccountNumber,
	}

	return res, nil
}

func (as accountService) DebitAccount(accountNumber string, amount int64, trx *gorm.DB) (*types.AccountResponse, error) {
	var lockedDB repositories.IAccountRepository

	if trx == nil {
		lockedDB = as.accountRepo.WithLock(&clause.Locking{
			Strength: "UPDATE",
		})
	} else {
		lockedDB = as.accountRepo.WithTx(trx).WithLock(&clause.Locking{
			Strength: "UPDATE",
		})
	}

	account, err := lockedDB.FindByAccountNumber(accountNumber)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			account = &models.Account{
				AccountNumber: accountNumber,
			}
			err = lockedDB.Create(account)
			if err != nil {
				return nil, errors.New("an error occurred when retrieving account, please try again")
			}
		} else {
			return nil, errors.New("an error occurred when retrieving account, please try again")
		}
	}
	if account.AvailableBalance < amount {
		return nil, errors.New("insufficient balance")
	}

	prevBal := account.AvailableBalance
	account.AvailableBalance -= amount
	err = lockedDB.Update(account)

	if err != nil {
		return nil, err
	}

	res := &types.AccountResponse{
		CurrentBalance:  account.AvailableBalance,
		PreviousBalance: prevBal,
		AccountNumber:   account.AccountNumber,
	}

	return res, nil
}

func (as accountService) LockFunds(userId string, accountNumber string, amount int64, duration int, trx *gorm.DB) (*types.SavingsResponse, error) {
	var lockedDB repositories.IAccountRepository

	if trx == nil {
		lockedDB = as.accountRepo.WithLock(&clause.Locking{
			Strength: "UPDATE",
		})
	} else {
		lockedDB = as.accountRepo.WithTx(trx).WithLock(&clause.Locking{
			Strength: "UPDATE",
		})
	}

	account, err := lockedDB.FindByAccountNumber(accountNumber)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			account = &models.Account{
				AccountNumber: accountNumber,
			}
			err = lockedDB.Create(account)
			if err != nil {
				return nil, errors.New("an error occurred when retrieving account, please try again")
			}
		} else {
			return nil, errors.New("an error occurred when retrieving account, please try again")
		}
	}
	if account.AvailableBalance < amount {
		return nil, errors.New("insufficient balance")
	}

	prevBal := account.AvailableBalance
	account.AvailableBalance -= amount
	err = lockedDB.Update(account)

	savings := models.Savings{
		UserId:        userId,
		AccountNumber: accountNumber,
		Amount:        amount,
		Duration:      duration,
	}
	err = as.savingsRepo.Create(&savings)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return &types.SavingsResponse{
		AccountNumber:   account.AccountNumber,
		CurrentBalance:  account.AvailableBalance,
		PreviousBalance: prevBal,
	}, nil
}
