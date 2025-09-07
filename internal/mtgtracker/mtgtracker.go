package mtgtracker

import (
	"encoding/json"
	"errors"
	"log"
	"mtgtracker/internal/middleware"
	"mtgtracker/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

type Storage interface {
	GeneratePresignedUploadURL(fileName string, contentType string) (string, error)
}
type Service struct {
	Repository *repository.Repository
	Storage    Storage
}

func NewService(repo *repository.Repository, storage Storage) *Service {
	return &Service{
		Repository: repo,
		Storage:    storage,
	}
}

func (s *Service) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /player/v1/signup", s.SignupPlayer)
	mux.HandleFunc("GET /player/v1/players", s.GetPlayers)
	mux.HandleFunc("GET /player/v1/players/{playerId}", s.GetPlayer)
	mux.HandleFunc("GET /player/v1/me", s.GetMyPlayer)
	mux.HandleFunc("POST /game/v1/games", s.CreateGame)
	mux.HandleFunc("GET /game/v1/games", s.GetGames)
	mux.HandleFunc("PUT /game/v1/games/{gameId}", s.UpdateGame)
	mux.HandleFunc("GET /game/v1/games/{gameId}", s.GetGame)
	mux.HandleFunc("DELETE /game/v1/games/{gameId}", s.DeleteGame)
	mux.HandleFunc("POST /game/v1/games/{gameId}/events", s.AddGameEvent)
	mux.HandleFunc("GET /ranking/v1/games/pending", s.GetPendingGames)
	mux.HandleFunc("PUT /ranking/v1/rankings/{rankingId}/accept", s.AcceptRanking)
	mux.HandleFunc("PUT /ranking/v1/rankings/{rankingId}/decline", s.DeclineRanking)
	mux.HandleFunc("DELETE /follow/v1/follows/{playerId}", s.DeleteFollow)
	mux.HandleFunc("GET /follow/v1/follows", s.GetMyFollows)
	mux.HandleFunc("GET /follow/v1/follows/{playerId}", s.GetPlayerFollows)
}

func (s *Service) GetPendingGames(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the context
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Call the repository to get the games with pending rankings
	games, err := s.Repository.GetPendingGames(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert games to DTOs
	result := make([]Game, 0, len(games))
	for _, game := range games {
		result = append(result, convertGameToDto(&game))
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}
func (s *Service) GetActiveGames(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the context
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Call the repository to get the games that are not finished
	games, err := s.Repository.GetActiveGames(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert games to DTOs
	result := make([]Game, 0, len(games))
	for _, game := range games {
		result = append(result, convertGameToDto(&game))
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}

}

func (s *Service) AcceptRanking(w http.ResponseWriter, r *http.Request) {
	rankingID := r.PathValue("rankingId")
	rankingIDInt, err := strconv.Atoi(rankingID)
	if err != nil {
		http.Error(w, "Invalid ranking ID", http.StatusBadRequest)
		return
	}

	// Call the repository to accept the ranking
	_, err = s.Repository.AcceptRanking(uint(rankingIDInt))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) DeclineRanking(w http.ResponseWriter, r *http.Request) {
	rankingID := r.PathValue("rankingId")
	rankingIDInt, err := strconv.Atoi(rankingID)
	if err != nil {
		http.Error(w, "Invalid ranking ID", http.StatusBadRequest)
		return
	}

	// Call the repository to decline the ranking
	err = s.Repository.DeclineRanking(uint(rankingIDInt))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) GetMyPlayer(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the context
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Call the repository to get the player by Firebase ID
	player, err := s.Repository.GetPlayerByFirebaseID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := convertPlayerToDto(player)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

func (s *Service) SignupPlayer(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var request SignupPlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the firebase user id from the context
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	//to do: get email from firebase
	email := middleware.GetUserEmail(r)
	if email == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	player, err := s.Repository.InsertPlayer(request.Name, email, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := convertPlayerToDto(player)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

func (s *Service) GetPlayers(w http.ResponseWriter, r *http.Request) {
	// Call the repository to get the players
	// use the search query parameters if needed
	search := r.URL.Query().Get("search")
	players, err := s.Repository.GetPlayers(search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]Player, 0, len(players))
	for _, player := range players {
		// Convert the player to a DTO
		result = append(result, convertPlayerToDto(&player))
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

func (s *Service) GetPlayer(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("playerId")
	// Call the repository to get the player by Firebase ID
	player, err := s.Repository.GetPlayerByFirebaseID(playerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := convertPlayerToDto(player)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
		return
	}
}
func (s *Service) CreateGame(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	userId := middleware.GetUserID(r)
	var request CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(request.Rankings) < 2 {
		http.Error(w, "At least two rankings are required", http.StatusBadRequest)
		return
	}
	// Call the repository to insert the game
	var rankings []repository.Ranking
	for _, rank := range request.Rankings {

		toAdd := repository.Ranking{
			PlayerID: rank.PlayerID,
			Position: 0,
			Status:   repository.StatusPending,
			Deck:     convertDeck(rank.Deck),
		}

		// we always accept our own ranking duh
		if rank.PlayerID != nil && userId == *rank.PlayerID {
			toAdd.Status = repository.StatusAccepted
		}
		rankings = append(rankings, toAdd)
	}
	game, err := s.Repository.InsertGame(request.Duration, request.Comments, request.Image, request.Date, request.Finished, rankings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := convertGameToDto(game)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

func (s *Service) AddGameEvent(w http.ResponseWriter, r *http.Request) {
	var req GameEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := r.PathValue("gameId")
	gameId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}
	// fetch the game and check if source and target rankings are valid
	_, err = s.Repository.GetGameWithEvents(uint(gameId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var uploadImgUrl string
	// If this is an image event, generate a presigned upload URL
	// add the current timestamp to the filename to avoid conflicts
	if req.EventType == repository.EventTypeImage {
		if req.EventImageName == nil {
			http.Error(w, "Event image name is required for image events", http.StatusBadRequest)
			return
		}
		uploadURL, err := s.Storage.GeneratePresignedUploadURL(*req.EventImageName, getImgContentType(*req.EventImageName))
		if err != nil {
			http.Error(w, "Error generating upload URL", http.StatusInternalServerError)
			return
		}
		uploadImgUrl = uploadURL
	}

	// Insert the event using the repository
	event, err := s.Repository.InsertGameEvent(uint(gameId), req.EventType, req.DamageDelta, req.TargetLifeTotalAfter, req.SourceRankingId, req.TargetRankingId, strings.Split(uploadImgUrl, "?")[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	eventDto := convertGameEvent(event, uploadImgUrl)

	err = json.NewEncoder(w).Encode(eventDto)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

func (s *Service) GetGame(w http.ResponseWriter, r *http.Request) {
	gameIdStr := r.PathValue("gameId")
	gameId, err := strconv.Atoi(gameIdStr)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}
	game, err := s.Repository.GetGameWithEvents(uint(gameId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// convert the game to a DTO
	result := convertGameToDto(game)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

func (s *Service) DeleteGame(w http.ResponseWriter, r *http.Request) {
	gameIdStr := r.PathValue("gameId")
	gameId, err := strconv.Atoi(gameIdStr)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	// Call the repository to delete the game
	err = s.Repository.DeleteGame(uint(gameId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) GetGames(w http.ResponseWriter, r *http.Request) {
	// Call the repository to get the games
	games, err := s.Repository.GetGames()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]Game, 0, len(games))
	for _, game := range games {
		// Convert the game to a DTO
		result = append(result, convertGameToDto(&game))
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

// UpdateGame updates an existing game by its ID.
// You may need to adjust the request struct and repository method as needed.
func (s *Service) UpdateGame(w http.ResponseWriter, r *http.Request) {
	gameIdStr := r.PathValue("gameId")
	gameId, err := strconv.Atoi(gameIdStr)
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	var request UpdateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// find the rankings for that game
	rankings := []repository.Ranking{}
	err = s.Repository.DB.Model(&repository.Ranking{}).Where("game_id = ?", gameId).Find(&rankings).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Validate and reorder rankings
	newRankings, err := validateAndReorderRankings(request.Rankings, rankings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the repository to update the game (implement UpdateGame in your repository)
	updatedGame, err := s.Repository.UpdateGame(uint(gameId), newRankings, request.Finished)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if updatedGame == nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	result := convertGameToDto(updatedGame)
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
	currentPlayer, err := s.Repository.GetPlayerByFirebaseID(userID)
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
	currentPlayer, err := s.Repository.GetPlayerByFirebaseID(userID)
	if err != nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	follows, err := s.Repository.GetFollows(currentPlayer.FirebaseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]Player, 0, len(follows))
	for _, player := range follows {
		result = append(result, convertPlayerToDto(&player))
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

	result := make([]Player, 0, len(follows))
	for _, player := range follows {
		result = append(result, convertPlayerToDto(&player))
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

// validateAndReorderRankings validates that the request rankings match existing rankings
// and returns them with sequential positions (1, 2, 3, etc.) based on request order
func validateAndReorderRankings(requestRankings []UpdateRanking, existingRankings []repository.Ranking) ([]repository.Ranking, error) {
	// Validate that request rankings match existing rankings count
	if len(requestRankings) != len(existingRankings) {
		return nil, errors.New("rankings count must match existing rankings")
	}

	// Create map of existing rankings by PlayerID string value for quick lookup
	existingMap := make(map[string]repository.Ranking)
	for _, ranking := range existingRankings {
		if ranking.PlayerID != nil {
			existingMap[*ranking.PlayerID] = ranking
		}
	}

	// Validate all request rankings have valid PlayerIDs and set sequential positions
	newRankings := make([]repository.Ranking, len(requestRankings))
	for i, reqRanking := range requestRankings {
		if reqRanking.PlayerID == nil {
			return nil, errors.New("player ID cannot be nil")
		}
		existing, exists := existingMap[*reqRanking.PlayerID]
		if !exists {
			return nil, errors.New("invalid player ID in rankings")
		}
		existing.Position = i + 1 // Set position to 1, 2, 3, etc.
		newRankings[i] = existing
	}

	return newRankings, nil
}
