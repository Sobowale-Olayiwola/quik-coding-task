package service

import (
	"context"
	"errors"
	"quik/domain"
	"quik/domain/mocks/inmemorydb"
	"quik/domain/mocks/repository"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCredit(t *testing.T) {
	as := assert.New(t)
	walletInMemoryDB := &inmemorydb.WalletInMemoryDBMock{}
	walletRepo := &repository.WalletRepositoryMock{}

	t.Run("happy path: Successfully credits a players balance", func(t *testing.T) {
		id := "6"
		amount := "5000"
		balance, _ := decimal.NewFromString("900")
		walletRepo.On("Get", context.Background(), id).Return(domain.Wallet{
			Balance:  balance,
			PlayerID: 1,
			ID:       6,
		}, nil).Once()
		walletRepo.On("Credit", context.Background(), mock.Anything).Return(nil)
		walletInMemoryDB.On("Delete", context.Background(), id).Return(nil).Once()
		service := NewWalletService(walletRepo, walletInMemoryDB)
		err := service.Credit(context.Background(), id, amount)
		as.NoError(err)
		walletRepo.AssertExpectations(t)
		walletInMemoryDB.AssertExpectations(t)
	})

	t.Run("input error: Negative amount ", func(t *testing.T) {
		id := "6"
		amount := "-5000"
		walletRepo.On("Get", context.Background(), id).Return(domain.Wallet{}, nil).Once()
		service := NewWalletService(walletRepo, walletInMemoryDB)
		err := service.Credit(context.Background(), id, amount)
		as.Error(err)
		walletRepo.AssertExpectations(t)
		walletInMemoryDB.AssertExpectations(t)
	})

	t.Run("input error: Invalid amount ", func(t *testing.T) {
		id := "6"
		amount := "-5000AERA"
		walletRepo.On("Get", context.Background(), id).Return(domain.Wallet{}, nil).Once()
		service := NewWalletService(walletRepo, walletInMemoryDB)
		err := service.Credit(context.Background(), id, amount)
		as.Error(err)
		walletRepo.AssertExpectations(t)
		walletInMemoryDB.AssertExpectations(t)
	})

	t.Run("system error: Database failed ", func(t *testing.T) {
		id := "6"
		amount := "5000"
		walletRepo.On("Get", context.Background(), id).Return(domain.Wallet{}, errors.New("something failed")).Once()
		service := NewWalletService(walletRepo, walletInMemoryDB)
		err := service.Credit(context.Background(), id, amount)
		as.Error(err)
		walletRepo.AssertExpectations(t)
		walletInMemoryDB.AssertExpectations(t)
	})
}

func TestDebit(t *testing.T) {
	as := assert.New(t)
	walletInMemoryDB := &inmemorydb.WalletInMemoryDBMock{}
	walletRepo := &repository.WalletRepositoryMock{}

	t.Run("happy path: Successfully debits a players balance", func(t *testing.T) {
		id := "6"
		amount := "900"
		balance, _ := decimal.NewFromString("5000")
		walletRepo.On("Get", context.Background(), id).Return(domain.Wallet{
			Balance:  balance,
			PlayerID: 1,
			ID:       6,
		}, nil).Once()
		walletRepo.On("Debit", context.Background(), mock.Anything).Return(nil).Once()
		walletInMemoryDB.On("Delete", context.Background(), id).Return(nil).Once()
		service := NewWalletService(walletRepo, walletInMemoryDB)
		err := service.Debit(context.Background(), id, amount)
		as.NoError(err)
		walletRepo.AssertExpectations(t)
		walletInMemoryDB.AssertExpectations(t)
	})

	t.Run("input error: Negative amount ", func(t *testing.T) {
		id := "6"
		amount := "-5000"
		walletRepo.On("Get", context.Background(), id, mock.Anything).Return(domain.Wallet{}, nil).Once()
		service := NewWalletService(walletRepo, walletInMemoryDB)
		err := service.Debit(context.Background(), id, amount)
		as.Error(err)
		walletRepo.AssertExpectations(t)
		walletInMemoryDB.AssertExpectations(t)
	})

	t.Run("input error: Invalid amount ", func(t *testing.T) {
		id := "6"
		amount := "-5000AERA"
		walletRepo.On("Get", context.Background(), id, mock.Anything).Return(domain.Wallet{}, nil).Once()
		service := NewWalletService(walletRepo, walletInMemoryDB)
		err := service.Debit(context.Background(), id, amount)
		as.Error(err)
		walletRepo.AssertExpectations(t)
		walletInMemoryDB.AssertExpectations(t)
	})

	t.Run("input error: Insufficient funds ", func(t *testing.T) {
		id := "6"
		amount := "5000"
		balance, _ := decimal.NewFromString("900")
		walletRepo.On("Get", context.Background(), id).Return(domain.Wallet{
			Balance:  balance,
			PlayerID: 1,
			ID:       6,
		}, nil).Once()
		service := NewWalletService(walletRepo, walletInMemoryDB)
		err := service.Debit(context.Background(), id, amount)
		as.Error(err)
		walletRepo.AssertExpectations(t)
		walletInMemoryDB.AssertExpectations(t)
	})

	t.Run("system error: Database failed ", func(t *testing.T) {
		id := "6"
		amount := "5000"
		walletRepo.On("Get", context.Background(), id).Return(domain.Wallet{}, errors.New("something failed")).Once()
		service := NewWalletService(walletRepo, walletInMemoryDB)
		err := service.Debit(context.Background(), id, amount)
		as.Error(err)
		walletRepo.AssertExpectations(t)
		walletInMemoryDB.AssertExpectations(t)
	})
}
