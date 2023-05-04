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
	CreateAccount(userId string) error
	CreditAccount(userId string, amount int64, trx *gorm.DB) (*types.AccountResponse, error)
	DebitAccount(userId string, amount int64, trx *gorm.DB) (*types.AccountResponse, error)
}

type accountService struct {
	userRepo    repositories.IUserRepository
	accountRepo repositories.IAccountRepository
}

func NewAccountService() IAccountService {
	return &accountService{
		userRepo:    repositories.NewUserRepo(),
		accountRepo: repositories.NewAccountRepo(),
	}
}

var (
	ErrAccountAlreadyExists = errors.New("an account already exists for this user")
)

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

func (as accountService) CreditAccount(userId string, amount int64, trx *gorm.DB) (*types.AccountResponse, error) {

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

	account, err := lockedDB.FindByUserId(userId)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			account = &models.Account{
				UserId: userId,
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

func (as accountService) DebitAccount(userId string, amount int64, trx *gorm.DB) (*types.AccountResponse, error) {
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

	account, err := lockedDB.FindByUserId(userId)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			account = &models.Account{
				UserId: userId,
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
