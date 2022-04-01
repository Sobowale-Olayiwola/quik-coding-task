package repository

import (
	"context"
	"errors"
	"quik/domain"

	"gorm.io/gorm"
)

type mysqlPlayerRepository struct {
	db *gorm.DB
}

func NewMySqlPlayerRepository(db *gorm.DB) domain.PlayerRepository {
	return &mysqlPlayerRepository{db}
}

func (m *mysqlPlayerRepository) Create(ctx context.Context, player *domain.Player) error {
	err := m.db.WithContext(ctx).Create(player).Error
	return err
}

func (m *mysqlPlayerRepository) Get(ctx context.Context, id string) (domain.Player, error) {
	var player domain.Player
	err := m.db.WithContext(ctx).Where("id = ?", id).First(&player).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return domain.Player{}, domain.ErrRecordNotFound
		default:
			return domain.Player{}, err
		}
	}
	return player, nil
}

func (m *mysqlPlayerRepository) Update(ctx context.Context, player *domain.Player, updatedPlayer domain.Player) error {
	err := m.db.WithContext(ctx).Model(player).Updates(updatedPlayer).Error
	return err
}

func (m *mysqlPlayerRepository) Delete(ctx context.Context, id string, player *domain.Player) error {
	err := m.db.WithContext(ctx).Where("id = ?", id).Delete(player).Error
	return err
}

func (m *mysqlPlayerRepository) FindByEmail(ctx context.Context, email string, player *domain.Player) error {
	err := m.db.WithContext(ctx).Where("email = ?", email).First(player).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return domain.ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}
