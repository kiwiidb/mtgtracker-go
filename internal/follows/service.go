package follows

import (
	"encoding/json"
	"log"
	"mtgtracker/internal/middleware"
	"mtgtracker/internal/mtgtracker"
	"mtgtracker/internal/repository"
	"net/http"
	"strings"
)

type Service struct {
	Repository    *Repository
	playerService playerService
}

type playerService interface {
	GetPlayerByFirebaseID(firebaseID string) (*repository.Player, error)
	ConvertPlayerToResponse(player *repository.Player) mtgtracker.Player
}

func NewService(repo *Repository, playerSvc playerService) *Service {
	return &Service{Repository: repo, playerService: playerSvc}
}

func (s *Service) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /follow/v1/follows/{playerId}", s.CreateFollow)
	mux.HandleFunc("DELETE /follow/v1/follows/{playerId}", s.DeleteFollow)
	mux.HandleFunc("GET /follow/v1/follows", s.GetMyFollows)
	mux.HandleFunc("GET /follow/v1/players/{playerId}/follows", s.GetPlayerFollows)
}

func (s *Service) CreateFollow(w http.ResponseWriter, r *http.Request) {
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

	follow, err := s.Repository.CreateFollow(currentPlayer.FirebaseID, targetPlayerID)
	if err != nil {
		if strings.Contains(err.Error(), "cannot follow yourself") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := struct {
		ID      uint              `json:"id"`
		Player1 mtgtracker.Player `json:"player1"`
		Player2 mtgtracker.Player `json:"player2"`
	}{
		ID:      follow.ID,
		Player1: s.playerService.ConvertPlayerToResponse(&follow.Player1),
		Player2: s.playerService.ConvertPlayerToResponse(&follow.Player2),
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

func (s *Service) DeleteFollow(w http.ResponseWriter, r *http.Request) {
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

	err = s.Repository.DeleteFollow(currentPlayer.FirebaseID, targetPlayerID)
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

func (s *Service) GetMyFollows(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the current user's player record
	currentPlayer, err := s.playerService.GetPlayerByFirebaseID(userID)
	if err != nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	follows, err := s.Repository.GetFollows(currentPlayer.FirebaseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]mtgtracker.Player, 0, len(follows))
	for _, player := range follows {
		result = append(result, s.playerService.ConvertPlayerToResponse(&player))
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

func (s *Service) GetPlayerFollows(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("playerId")

	follows, err := s.Repository.GetFollows(playerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]mtgtracker.Player, 0, len(follows))
	for _, player := range follows {
		result = append(result, s.playerService.ConvertPlayerToResponse(&player))
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}
