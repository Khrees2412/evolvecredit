package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	Base
	Id               string `json:"id"`
	UserId           string `json:"user_id" gorm:"foreignKey"`
	AccountNumber    string `json:"account_number" gorm:"unique"`
	LedgerBalance    int64  `json:"ledger_balance" gorm:"not null;default:0"`
	AvailableBalance int64  `json:"available_balance" gorm:"not null;default:0"`
	TotalLocked      int64  `json:"total_locked" gorm:"not null;default:0"`
	IsActive         *bool  `json:"is_active" gorm:"not null;default:true"`
}

type Savings struct {
	Base
	Id       string  `json:"id"`
	UserId   string  `json:"user_id"`
	Amount   float64 `json:"amount"`
	Duration int     `json:"duration"` // savings lock duration in days
}

func (u *Account) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id = uuid.NewString()
	return nil
}

func (u *Savings) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id = uuid.NewString()
	return nil
}
