package service

import (
	"context"
	"quik/domain"
	"sync"

	"github.com/shopspring/decimal"
)

type walletService struct {
	walletRepository domain.WalletRepository
	walletInMemoryDB domain.WalletInMemoryDB
}

func NewWalletService(r domain.WalletRepository, i domain.WalletInMemoryDB) domain.WalletService {
	return &walletService{
		walletRepository: r,
		walletInMemoryDB: i,
	}
}

func (w *walletService) Create(ctx context.Context, wallet *domain.Wallet) error {
	err := w.walletRepository.Create(ctx, wallet)
	return err
}

func (w *walletService) Get(ctx context.Context, id string) (domain.Wallet, error) {
	wallet, err := w.walletInMemoryDB.Get(ctx, id)
	if err == domain.ErrKeyNotFound {
		wallet, err = w.walletRepository.Get(ctx, id)
		if err != nil {
			return domain.Wallet{}, err
		}
		w.walletInMemoryDB.Set(ctx, id, &wallet)
	}
	return wallet, nil
}

func (w *walletService) Credit(ctx context.Context, id, amount string) error {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	wallet, err := w.walletRepository.Get(ctx, id)
	if err != nil {
		return err
	}
	creditAmount, err := decimal.NewFromString(amount)
	if creditAmount.IsNegative() || err != nil {
		return domain.ErrInvalidAmount
	}
	wallet.Balance = wallet.Balance.Add(creditAmount)
	err = w.walletRepository.Credit(ctx, &wallet)
	w.walletInMemoryDB.Delete(ctx, id)
	mutex.Unlock()
	return err
}

func (w *walletService) Debit(ctx context.Context, id, amount string) error {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	wallet, err := w.walletRepository.Get(ctx, id)
	if err != nil {
		return err
	}
	debitAmount, err := decimal.NewFromString(amount)
	if debitAmount.IsNegative() || err != nil {
		return domain.ErrInvalidAmount
	}
	wallet.Balance = wallet.Balance.Sub(debitAmount)
	if wallet.Balance.IsNegative() {
		return domain.ErrInsufficientFunds
	}
	err = w.walletRepository.Debit(ctx, &wallet)
	w.walletInMemoryDB.Delete(ctx, id)
	mutex.Unlock()
	return err
}
