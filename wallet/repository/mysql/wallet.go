package mysql

import (
	"context"
	"errors"
	"quik/domain"

	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type mysqlWalletRepository struct {
	db *gorm.DB
}

func NewMySqlWalletRepository(db *gorm.DB) domain.WalletRepository {
	return &mysqlWalletRepository{db: db}
}

func (w *mysqlWalletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	err := w.db.WithContext(ctx).Create(wallet).Error
	return err
}

func (w *mysqlWalletRepository) Get(ctx context.Context, id string) (domain.Wallet, error) {
	var wallet domain.Wallet
	err := w.db.WithContext(ctx).Where("id = ?", id).First(&wallet).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return domain.Wallet{}, ErrRecordNotFound
		default:
			return domain.Wallet{}, err
		}
	}
	return wallet, nil
}

func (w *mysqlWalletRepository) Credit(ctx context.Context, wallet *domain.Wallet) error {
	err := w.db.WithContext(ctx).Save(wallet).Error
	return err
}

func (w *mysqlWalletRepository) Debit(ctx context.Context, wallet *domain.Wallet) error {
	err := w.db.WithContext(ctx).Save(wallet).Error
	return err
}
