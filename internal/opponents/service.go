package opponents

import (
	"encoding/json"
	"log"
	"mtgtracker/internal/core"
	"mtgtracker/internal/middleware"
	"mtgtracker/internal/pagination"
	"net/http"
	"strings"
)

type Service struct {
	Repository    *Repository
	playerService playerService
}

type playerService interface {
	GetPlayerByFirebaseID(firebaseID string) (*core.Player, error)
	ConvertPlayerToResponse(player *core.Player) core.PlayerResponse
	GetGameByID(gameID uint) (*core.Game, error)
}

func NewService(repo *Repository, playerSvc playerService) *Service {
	return &Service{Repository: repo, playerService: playerSvc}
}

func (s *Service) RegisterRoutes(mux *http.ServeMux) {
	//mux.HandleFunc("POST /opponent/v1/opponents/{playerId}", s.CreateOpponent)
	//mux.HandleFunc("DELETE /opponent/v1/opponents/{playerId}", s.DeleteOpponent)
	mux.HandleFunc("GET /opponent/v1/opponents", s.GetMyOpponents)
	mux.HandleFunc("GET /opponent/v1/players/{playerId}/opponents", s.GetPlayerOpponents)
}

func (s *Service) CreateOpponent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetPlayerID := r.PathValue("playerId")

	// Get the current user's player record
	currentPlayer, err := s.playerService.GetPlayerByFirebaseID(userID)
	if err != nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	opponent, err := s.Repository.CreateOpponent(currentPlayer.FirebaseID, targetPlayerID)
	if err != nil {
		if strings.Contains(err.Error(), "cannot opponent yourself") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := struct {
		ID      uint                `json:"id"`
		Player1 core.PlayerResponse `json:"player1"`
		Player2 core.PlayerResponse `json:"player2"`
	}{
		ID:      opponent.ID,
		Player1: s.playerService.ConvertPlayerToResponse(&opponent.Player1),
		Player2: s.playerService.ConvertPlayerToResponse(&opponent.Player2),
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

func (s *Service) DeleteOpponent(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetPlayerID := r.PathValue("playerId")

	// Get the current user's player record
	currentPlayer, err := s.playerService.GetPlayerByFirebaseID(userID)
	if err != nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	err = s.Repository.DeleteOpponent(currentPlayer.FirebaseID, targetPlayerID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) GetMyOpponents(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	p := pagination.ParsePagination(r)

	// Get the current user's player record
	currentPlayer, err := s.playerService.GetPlayerByFirebaseID(userID)
	if err != nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	opponents, total, err := s.Repository.GetOpponents(currentPlayer.FirebaseID, p.PerPage, p.Offset())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items := make([]core.PlayerOpponentWithCount, 0, len(opponents))
	for _, opponent := range opponents {
		items = append(items, core.PlayerOpponentWithCount{
			Player: s.playerService.ConvertPlayerToResponse(&opponent.Player),
			Count:  opponent.GameCount,
		})
	}

	result := pagination.PaginatedResult[core.PlayerOpponentWithCount]{
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

func (s *Service) GetPlayerOpponents(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("playerId")
	p := pagination.ParsePagination(r)

	opponents, total, err := s.Repository.GetOpponents(playerID, p.PerPage, p.Offset())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items := make([]core.PlayerOpponentWithCount, 0, len(opponents))
	for _, opponent := range opponents {
		items = append(items, core.PlayerOpponentWithCount{
			Player: s.playerService.ConvertPlayerToResponse(&opponent.Player),
			Count:  opponent.GameCount,
		})
	}

	result := pagination.PaginatedResult[core.PlayerOpponentWithCount]{
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
