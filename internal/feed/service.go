package feed

import (
	"encoding/json"
	"log"
	"mtgtracker/internal/core"
	"mtgtracker/internal/middleware"
	"mtgtracker/internal/opponents"
	"mtgtracker/internal/pagination"
	"net/http"
)

type OpponentRepository interface {
	GetOpponents(playerID string, limit, offset int) ([]opponents.OpponentWithCount, int64, error)
}

type GameRepository interface {
	SearchGamesWithFilters(filter core.GameFilter, limit, offset int) ([]core.Game, int64, error)
}

type GameConverter interface {
	ConvertGameToDto(game *core.Game, includeEvents bool) core.GameResponse
}

type Service struct {
	opponentRepo  OpponentRepository
	gameRepo      GameRepository
	gameConverter GameConverter
}

func NewService(opponentRepo OpponentRepository, gameRepo GameRepository, gameConverter GameConverter) *Service {
	return &Service{
		opponentRepo:  opponentRepo,
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

	// Get all opponents for the player (we need all of them, not paginated)
	opponents, _, err := s.opponentRepo.GetOpponents(userID, 1000, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract player IDs from opponents
	playerIDs := make([]string, 0, len(opponents))
	for _, opponent := range opponents {
		playerIDs = append(playerIDs, opponent.Player.FirebaseID)
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
