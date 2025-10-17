package statistics

import (
	"encoding/json"
	"log"
	"mtgtracker/internal/middleware"
	"mtgtracker/internal/pagination"
	"net/http"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// RegisterRoutes registers HTTP endpoints for statistics
func (s *Service) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /statistics/v1/players", s.GetAllLatestPlayerStats)
	mux.HandleFunc("GET /statistics/v1/players/{playerId}", s.GetLatestPlayerStats)
	mux.HandleFunc("GET /statistics/v1/players/{playerId}/timeseries", s.GetPlayerStatsTimeSeries)
	mux.HandleFunc("GET /statistics/v1/me", s.GetMyLatestStats)
	mux.HandleFunc("GET /statistics/v1/me/timeseries", s.GetMyStatsTimeSeries)
}

// GetLatestPlayerStats retrieves the most recent statistics for a specific player
func (s *Service) GetLatestPlayerStats(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("playerId")
	if playerID == "" {
		http.Error(w, "Player ID is required", http.StatusBadRequest)
		return
	}

	stats, err := s.repo.GetLatestPlayerStats(playerID)
	if err != nil {
		http.Error(w, "Statistics not found", http.StatusNotFound)
		return
	}

	response := stats.ToResponse()
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

// GetPlayerStatsTimeSeries retrieves time series statistics for a specific player
func (s *Service) GetPlayerStatsTimeSeries(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("playerId")
	if playerID == "" {
		http.Error(w, "Player ID is required", http.StatusBadRequest)
		return
	}

	p := pagination.ParsePagination(r)

	stats, total, err := s.repo.GetPlayerStatsTimeSeries(playerID, p.PerPage, p.Offset())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]PlayerStatsResponse, len(stats))
	for i, stat := range stats {
		responses[i] = stat.ToResponse()
	}

	result := pagination.PaginatedResult[PlayerStatsResponse]{
		Items:      responses,
		TotalCount: total,
		Page:       p.Page,
		PerPage:    p.PerPage,
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

// GetMyLatestStats retrieves the most recent statistics for the authenticated user
func (s *Service) GetMyLatestStats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	stats, err := s.repo.GetLatestPlayerStats(userID)
	if err != nil {
		http.Error(w, "Statistics not found", http.StatusNotFound)
		return
	}

	response := stats.ToResponse()
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

// GetMyStatsTimeSeries retrieves time series statistics for the authenticated user
func (s *Service) GetMyStatsTimeSeries(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	p := pagination.ParsePagination(r)

	stats, total, err := s.repo.GetPlayerStatsTimeSeries(userID, p.PerPage, p.Offset())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]PlayerStatsResponse, len(stats))
	for i, stat := range stats {
		responses[i] = stat.ToResponse()
	}

	result := pagination.PaginatedResult[PlayerStatsResponse]{
		Items:      responses,
		TotalCount: total,
		Page:       p.Page,
		PerPage:    p.PerPage,
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

// GetAllLatestPlayerStats retrieves the most recent statistics for all players with pagination
func (s *Service) GetAllLatestPlayerStats(w http.ResponseWriter, r *http.Request) {
	p := pagination.ParsePagination(r)

	stats, total, err := s.repo.GetAllLatestPlayerStats(p.PerPage, p.Offset())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]PlayerStatsResponse, len(stats))
	for i, stat := range stats {
		responses[i] = stat.ToResponse()
	}

	result := pagination.PaginatedResult[PlayerStatsResponse]{
		Items:      responses,
		TotalCount: total,
		Page:       p.Page,
		PerPage:    p.PerPage,
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}
