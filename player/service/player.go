package service

import (
	"context"
	"errors"
	"quik/domain"
	"strings"
)

type playerService struct {
	playerRepository domain.PlayerRepository
}

func NewPlayerService(dp domain.PlayerRepository) domain.PlayerService {
	return &playerService{dp}
}

func (p *playerService) FindByEmail(ctx context.Context, email string, player *domain.Player) error {
	err := p.playerRepository.FindByEmail(ctx, email, player)
	return err
}

func (p *playerService) Create(ctx context.Context, player *domain.Player) error {

	var userExists domain.Player
	err := p.playerRepository.FindByEmail(ctx, player.Email, &userExists)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			err = p.playerRepository.Create(ctx, player)
			return err
		default:
			return err
		}
	}
	if strings.EqualFold(userExists.Email, player.Email) {
		return domain.ErrDuplicateRecord
	}
	err = p.playerRepository.Create(ctx, player)
	return err
}

func (p *playerService) Get(ctx context.Context, id string) (domain.Player, error) {
	player, err := p.playerRepository.Get(ctx, id)
	return player, err
}

func (p *playerService) Update(ctx context.Context, id string, player *domain.Player, updatedPlayer domain.Player) error {

	err := p.playerRepository.Update(ctx, player, updatedPlayer)
	return err
}

func (p *playerService) Delete(ctx context.Context, id string, player *domain.Player) error {
	err := p.playerRepository.Delete(ctx, id, player)
	return err
}
