package models

import (
	"github.com/google/uuid"
	"github.com/khrees2412/evolvecredit/types"
	"gorm.io/gorm"
)

type Transaction struct {
	Base
	Id     string                  `json:"id"`
	UserId string                  `json:"user_id"`
	Amount float64                 `json:"amount"`
	Type   types.TransactionType   `json:"type"`
	Entry  types.TransactionEntry  `json:"entry"`
	Status types.TransactionStatus `json:"status"`
}

func (u *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id = uuid.NewString()
	return nil
}
