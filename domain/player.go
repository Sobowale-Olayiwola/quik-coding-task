package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrDuplicateRecord = errors.New("duplicate record")
)

type Player struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type PlayerService interface {
	Create(ctx context.Context, player *Player) error
	Get(ctx context.Context, id string) (Player, error)
	Update(ctx context.Context, id string, player *Player, updatedPlayer Player) error
	Delete(ctx context.Context, id string, player *Player) error
	FindByEmail(ctx context.Context, email string, player *Player) error
}

type PlayerRepository interface {
	Create(ctx context.Context, player *Player) error
	Update(ctx context.Context, player *Player, updatedPlayer Player) error
	Get(ctx context.Context, id string) (Player, error)
	Delete(ctx context.Context, id string, player *Player) error
	FindByEmail(ctx context.Context, email string, player *Player) error
}
