package repository

import (
	"context"
	"quik/domain"

	"github.com/stretchr/testify/mock"
)

type WalletRepositoryMock struct {
	mock.Mock
}

func (w *WalletRepositoryMock) Create(ctx context.Context, wallet *domain.Wallet) error {
	output := w.Mock.Called(ctx, wallet)
	err := output.Error(0)
	return err
}

func (w *WalletRepositoryMock) Get(ctx context.Context, id string) (domain.Wallet, error) {
	output := w.Mock.Called(ctx, id)
	wallet := output.Get(0)
	err := output.Error(1)
	return wallet.(domain.Wallet), err
}

func (w *WalletRepositoryMock) Credit(ctx context.Context, wallet *domain.Wallet) error {
	output := w.Mock.Called(ctx, wallet)
	err := output.Error(0)
	return err
}

func (w *WalletRepositoryMock) Debit(ctx context.Context, wallet *domain.Wallet) error {
	output := w.Mock.Called(ctx, wallet)
	err := output.Error(0)
	return err
}
