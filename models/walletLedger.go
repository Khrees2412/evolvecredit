package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/khrees2412/evolvecredit/types"
)

type WalletLedger struct {
	Id              string                 `json:"id"`
	UserId          string                 `json:"user_id"`
	TransactionId   string                 `json:"transaction_id"`
	AccountNumber   string                 `json:"account_number"`
	PreviousBalance int64                  `json:"previous_balance"`
	CurrentBalance  int64                  `json:"current_balance"`
	Entry           types.TransactionEntry `json:"entry"`
}

func (b *WalletLedger) BeforeCreate(tx *gorm.DB) (err error) {
	b.Id = uuid.NewString()
	return nil
}
