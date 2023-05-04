package repositories

import (
	"github.com/khrees2412/evolvecredit/database"
	"github.com/khrees2412/evolvecredit/models"
	"gorm.io/gorm"
)

type IWalletLedgerRepository interface {
	Create(walletLedger *models.WalletLedger) error
	Update(walletLedger *models.WalletLedger) error
	WithTx(tx *gorm.DB) IWalletLedgerRepository
}

type walletLedgerRepo struct {
	db *gorm.DB
}

// NewWalletLedgerRepo will instantiate WalletLedger Repository
func NewWalletLedgerRepo() IWalletLedgerRepository {
	return &walletLedgerRepo{
		db: database.DB(),
	}
}

func (t *walletLedgerRepo) Create(walletLedger *models.WalletLedger) error {
	return t.db.Create(walletLedger).Error
}

func (t *walletLedgerRepo) Update(walletLedger *models.WalletLedger) error {
	return t.db.Save(walletLedger).Error
}

func (t *walletLedgerRepo) WithTx(tx *gorm.DB) IWalletLedgerRepository {
	return &walletLedgerRepo{db: tx}
}
