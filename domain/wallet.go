package domain

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

var (
	ErrKeyNotFound       = errors.New("key not found")
	ErrInsufficientFunds = errors.New("insufficient fund")
	ErrInvalidAmount     = errors.New("invalid amount")
)

type Wallet struct {
	ID        int             `json:"id"`
	PlayerID  int             `json:"playerId"`
	Balance   decimal.Decimal `json:"balance"`
	UpdatedAt time.Time       `json:"updated_at"`
	CreatedAt time.Time       `json:"created_at"`
}

type WalletService interface {
	Create(ctx context.Context, w *Wallet) error
	Get(ctx context.Context, id string) (Wallet, error)
	Credit(ctx context.Context, id, amount string) error
	Debit(ctx context.Context, id, amount string) error
}

type WalletRepository interface {
	Create(ctx context.Context, w *Wallet) error
	Get(ctx context.Context, id string) (Wallet, error)
	Credit(ctx context.Context, w *Wallet) error
	Debit(ctx context.Context, w *Wallet) error
}

type WalletInMemoryDB interface {
	Get(ctx context.Context, id string) (Wallet, error)
	Set(ctx context.Context, id string, w *Wallet) error
	Delete(ctx context.Context, id string) error
}
