package inmemorydb

import (
	"context"
	"quik/domain"

	"github.com/stretchr/testify/mock"
)

type WalletInMemoryDBMock struct {
	mock.Mock
}

func (w *WalletInMemoryDBMock) Get(ctx context.Context, id string) (domain.Wallet, error) {
	output := w.Mock.Called(ctx, id)
	wallet := output.Get(0)
	err := output.Error(1)
	return wallet.(domain.Wallet), err
}

func (w *WalletInMemoryDBMock) Set(ctx context.Context, id string, wallet *domain.Wallet) error {
	output := w.Mock.Called(ctx, id, wallet)
	err := output.Error(0)
	return err
}

func (w *WalletInMemoryDBMock) Delete(ctx context.Context, id string) error {
	output := w.Mock.Called(ctx, id)
	err := output.Error(0)
	return err
}
