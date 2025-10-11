package notification

import (
	"encoding/json"
	"log"
	"mtgtracker/internal/core"
	"mtgtracker/internal/middleware"
	"mtgtracker/internal/pagination"
	"net/http"
	"strconv"
	"strings"
)

type Service struct {
	Repository  *Repository
	coreService CoreService
}

func NewService(repo *Repository, coreService CoreService) *Service {
	return &Service{Repository: repo, coreService: coreService}
}

type CoreService interface {
	ConvertPlayerToResponse(player *core.Player) core.PlayerResponse
	ConvertGameToDto(game *core.Game, includePlayers bool) core.GameResponse
	GetGameByID(gameID uint) (*core.Game, error)
	GetPlayerByFirebaseID(firebaseID string) (*core.Player, error)
}

func (s *Service) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /notification/v1/notifications", s.GetNotifications)
	mux.HandleFunc("PUT /notification/v1/notifications/{notificationId}/read", s.MarkNotificationAsRead)
}

func (s *Service) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	p := pagination.ParsePagination(r)

	// Get the read filter from query params (optional)
	readParam := r.URL.Query().Get("read")
	var readFilter *bool
	if readParam != "" {
		switch readParam {
		case "true":
			trueVal := true
			readFilter = &trueVal
		case "false":
			falseVal := false
			readFilter = &falseVal
		}
	}

	notifications, total, err := s.Repository.GetNotifications(userID, readFilter, p.PerPage, p.Offset())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items := make([]NotificationResponse, 0, len(notifications))
	for _, notification := range notifications {
		items = append(items, s.convertNotificationToDto(&notification))
	}

	result := pagination.PaginatedResult[NotificationResponse]{
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

func (s *Service) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notificationIDStr := r.PathValue("notificationId")
	notificationID, err := strconv.Atoi(notificationIDStr)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	err = s.Repository.MarkNotificationAsRead(uint(notificationID), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "access denied") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
