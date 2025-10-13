package feed

import (
	"encoding/json"
	"log"
	"mtgtracker/internal/core"
	"mtgtracker/internal/follows"
	"mtgtracker/internal/middleware"
	"mtgtracker/internal/pagination"
	"net/http"
)

type FollowRepository interface {
	GetFollows(playerID string, limit, offset int) ([]follows.FollowWithCount, int64, error)
}

type GameRepository interface {
	SearchGamesWithFilters(filter core.GameFilter, limit, offset int) ([]core.Game, int64, error)
}

type GameConverter interface {
	ConvertGameToDto(game *core.Game, includeEvents bool) core.GameResponse
}

type Service struct {
	followRepo    FollowRepository
	gameRepo      GameRepository
	gameConverter GameConverter
}

func NewService(followRepo FollowRepository, gameRepo GameRepository, gameConverter GameConverter) *Service {
	return &Service{
		followRepo:    followRepo,
		gameRepo:      gameRepo,
		gameConverter: gameConverter,
	}
}

func (s *Service) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /feed/v1/items", s.GetFeedItems)
}

func (s *Service) GetFeedItems(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	p := pagination.ParsePagination(r)

	// Get all follows for the player (we need all of them, not paginated)
	follows, _, err := s.followRepo.GetFollows(userID, 1000, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract player IDs from follows
	playerIDs := make([]string, 0, len(follows))
	for _, follow := range follows {
		playerIDs = append(playerIDs, follow.Player.FirebaseID)
	}

	// Search games with those player IDs
	filter := core.GameFilter{PlayerIDs: playerIDs}
	games, total, err := s.gameRepo.SearchGamesWithFilters(filter, p.PerPage, p.Offset())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to DTO
	items := make([]core.GameResponse, 0, len(games))
	for _, game := range games {
		items = append(items, s.gameConverter.ConvertGameToDto(&game, true))
	}

	result := pagination.PaginatedResult[core.GameResponse]{
		Items:      items,
		TotalCount: total,
		Page:       p.Page,
		PerPage:    p.PerPage,
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}
