package feed

import (
	"mtgtracker/internal/core"
)

type FollowRepository interface {
	GetFollows(playerID string) ([]string, error)
}

type GameRepository interface {
	SearchGames(playerIDs []string, limit, offset int) ([]core.Game, int64, error)
}

type Service struct {
	followRepo FollowRepository
	gameRepo   GameRepository
}

func NewService(followRepo FollowRepository, gameRepo GameRepository) *Service {
	return &Service{
		followRepo: followRepo,
		gameRepo:   gameRepo,
	}
}
