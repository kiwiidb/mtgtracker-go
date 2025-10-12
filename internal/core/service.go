package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mtgtracker/internal/events"
	"mtgtracker/internal/middleware"
	"mtgtracker/internal/pagination"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Storage interface {
	GeneratePresignedUploadURL(fileName string, contentType string) (string, error)
}

type EventBus interface {
	Publish(event events.Event)
}

type Service struct {
	Repository *Repository
	Storage    Storage
	eventBus   EventBus
}

func NewService(repo *Repository, storage Storage, eventBus EventBus) *Service {
	return &Service{
		Repository: repo,
		Storage:    storage,
		eventBus:   eventBus,
	}
}

const (
	profileImagePrefix = "profile-images/"
	eventImagePrefix   = "event-images/"
)

func (s *Service) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /player/v1/signup", s.SignupPlayer)
	mux.HandleFunc("GET /player/v1/players", s.GetPlayers)
	mux.HandleFunc("GET /player/v1/players/{playerId}", s.GetPlayer)
	mux.HandleFunc("GET /player/v1/me", s.GetMyPlayer)
	mux.HandleFunc("PUT /player/v1/me", s.UpdateMyPlayer)
	mux.HandleFunc("POST /player/v1/profile-image/upload-url", s.GetProfileImageUploadURL)
	mux.HandleFunc("GET /player/v1/players/{playerId}/decks", s.GetPlayerDecks)
	mux.HandleFunc("GET /player/v1/players/{playerId}/games", s.GetPlayerGames)
	mux.HandleFunc("POST /deck/v1/decks", s.CreateDeck)
	mux.HandleFunc("POST /game/v1/games", s.CreateGame)
	mux.HandleFunc("GET /game/v1/games", s.GetGames)
	mux.HandleFunc("GET /game/v1/games/active", s.GetActiveGame)
	mux.HandleFunc("PUT /game/v1/games/{gameId}", s.UpdateGame)
	mux.HandleFunc("GET /game/v1/games/{gameId}", s.GetGame)
	mux.HandleFunc("DELETE /game/v1/games/{gameId}", s.DeleteGame)
	mux.HandleFunc("DELETE /ranking/v1/rankings/{rankingId}", s.DeleteRanking)
	mux.HandleFunc("POST /game/v1/games/{gameId}/events", s.AddGameEvent)
}

func (s *Service) GetPlayerByFirebaseID(firebaseID string) (*Player, error) {
	return s.Repository.GetPlayerByFirebaseID(firebaseID)
}

func (s *Service) GetGameByID(gameID uint) (*Game, error) {
	return s.Repository.GetGameWithEvents(gameID)
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

	result := s.ConvertPlayerToResponse(player)
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

	result := s.ConvertPlayerToResponse(player)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

func (s *Service) GetPlayers(w http.ResponseWriter, r *http.Request) {
	p := pagination.ParsePagination(r)

	// Get search query parameter
	search := r.URL.Query().Get("search")

	// Call the repository to get the players
	players, total, err := s.Repository.GetPlayers(search, p.PerPage, p.Offset())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items := make([]PlayerResponse, 0, len(players))
	for _, player := range players {
		items = append(items, s.ConvertPlayerToResponse(&player))
	}

	result := pagination.PaginatedResult[PlayerResponse]{
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

func (s *Service) GetPlayer(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("playerId")
	// Call the repository to get the player by Firebase ID
	player, err := s.Repository.GetPlayerByFirebaseID(playerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result := s.ConvertPlayerToResponse(player)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
		return
	}
}
func (s *Service) CreateGame(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the context
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// find the current user
	user, err := s.Repository.GetPlayerByFirebaseID(userID)
	if err != nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}
	var request CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(request.Rankings) < 2 {
		http.Error(w, "At least two rankings are required", http.StatusBadRequest)
		return
	}

	// Validate rankings have either deck_id or inline deck
	for i, rank := range request.Rankings {
		if rank.DeckID == nil && rank.Deck == nil {
			http.Error(w, "Each ranking must have either deck_id or deck provided", http.StatusBadRequest)
			return
		}
		if rank.DeckID != nil && rank.Deck != nil {
			http.Error(w, "Each ranking must have either deck_id OR deck, not both", http.StatusBadRequest)
			return
		}
		// If inline deck is provided, validate required fields
		if rank.Deck != nil && rank.Deck.Commander == "" {
			http.Error(w, fmt.Sprintf("Ranking %d: deck.commander is required", i), http.StatusBadRequest)
			return
		}
	}

	// Call the repository to insert the game
	var rankings []Ranking
	for _, rank := range request.Rankings {
		toAdd := Ranking{
			PlayerID: rank.PlayerID,
			Position: 0,
			DeckID:   rank.DeckID, // Optional deck reference
		}

		// If inline deck is provided (and no deck_id), use embedded deck
		if rank.Deck != nil && rank.DeckID == nil {
			toAdd.DeckEmbedded = convertSimpleDeck(*rank.Deck)
		}

		rankings = append(rankings, toAdd)
	}
	game, err := s.Repository.InsertGame(user, request.Duration, request.Comments, request.Image, request.Date, request.Finished, rankings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Publish game created event
	rankingIDs := make([]uint, len(game.Rankings))
	for i, ranking := range game.Rankings {
		rankingIDs[i] = ranking.ID
	}
	s.eventBus.Publish(events.GameCreatedEvent{
		GameID:     game.ID,
		CreatorID:  user.FirebaseID,
		RankingIDs: rankingIDs,
		Date:       time.Now(),
	})

	result := s.ConvertGameToDto(game, false)
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
	if req.EventType == EventTypeImage {
		if req.EventImageName == nil {
			http.Error(w, "Event image name is required for image events", http.StatusBadRequest)
			return
		}
		uploadURL, err := s.Storage.GeneratePresignedUploadURL(eventImagePrefix+*req.EventImageName, getImgContentType(*req.EventImageName))
		if err != nil {
			http.Error(w, "Error generating upload URL", http.StatusInternalServerError)
			return
		}
		uploadImgUrl = uploadURL
	}

	// Insert the event using the repository
	event, err := s.Repository.InsertGameEvent(
		uint(gameId), req.EventType,
		req.DamageDelta, req.TargetLifeTotalAfter,
		req.SourceRankingId, req.TargetRankingId,
		strings.Split(uploadImgUrl, "?")[0],
		req.Comment,
	)
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
	result := s.ConvertGameToDto(game, true)
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

	// Fetch game data before deletion to include in event
	game, err := s.Repository.GetGameWithEvents(uint(gameId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Call the repository to delete the game
	err = s.Repository.DeleteGame(uint(gameId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Publish game deleted event with player/ranking info
	rankingIDs := make([]uint, len(game.Rankings))
	playerIDs := make([]string, 0, len(game.Rankings))
	for i, ranking := range game.Rankings {
		rankingIDs[i] = ranking.ID
		if ranking.PlayerID != nil {
			playerIDs = append(playerIDs, *ranking.PlayerID)
		}
	}

	s.eventBus.Publish(events.GameDeletedEvent{
		GameID:     uint(gameId),
		RankingIDs: rankingIDs,
		PlayerIDs:  playerIDs,
		Date:       time.Now(),
	})

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) GetGames(w http.ResponseWriter, r *http.Request) {
	p := pagination.ParsePagination(r)

	// Call the repository to get the games
	games, total, err := s.Repository.GetGames(p.PerPage, p.Offset())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items := make([]GameResponse, 0, len(games))
	for _, game := range games {
		items = append(items, s.ConvertGameToDto(&game, true))
	}

	result := pagination.PaginatedResult[GameResponse]{
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

func (s *Service) GetActiveGame(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	game, err := s.Repository.GetActiveGameForPlayer(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if game == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	result := s.ConvertGameToDto(game, true)
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
	rankings := []Ranking{}
	err = s.Repository.DB.Model(&Ranking{}).Where("game_id = ?", gameId).Find(&rankings).Error
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

	// Publish game finished event if game was just finished
	if request.Finished != nil && *request.Finished {
		rankingIDs := make([]uint, len(updatedGame.Rankings))
		for i, ranking := range updatedGame.Rankings {
			rankingIDs[i] = ranking.ID
		}
		s.eventBus.Publish(events.GameFinishedEvent{
			GameID:     updatedGame.ID,
			RankingIDs: rankingIDs,
			Date:       time.Now(),
		})
	}

	result := s.ConvertGameToDto(updatedGame, false)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

// validateAndReorderRankings validates that the request rankings match existing rankings
// and returns them with sequential positions (1, 2, 3, etc.) based on request order
func validateAndReorderRankings(requestRankings []UpdateRanking, existingRankings []Ranking) ([]Ranking, error) {
	// Validate that request rankings match existing rankings count
	if len(requestRankings) != len(existingRankings) {
		return nil, errors.New("rankings count must match existing rankings")
	}

	// Create map of existing rankings by RankingID string value for quick lookup
	existingMap := make(map[uint]Ranking)
	for _, ranking := range existingRankings {
		existingMap[ranking.ID] = ranking
	}

	newRankings := make([]Ranking, len(requestRankings))
	for i, reqRanking := range requestRankings {
		existing, exists := existingMap[reqRanking.RankingID]
		if !exists {
			return nil, errors.New("invalid ranking ID in rankings")
		}
		existing.Position = i + 1 // Set position to 1, 2, 3, etc.
		newRankings[i] = existing
	}

	return newRankings, nil
}

func (s *Service) GetProfileImageUploadURL(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request struct {
		FileName string `json:"file_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if request.FileName == "" {
		http.Error(w, "file_name is required", http.StatusBadRequest)
		return
	}

	uploadURL, err := s.Storage.GeneratePresignedUploadURL(profileImagePrefix+request.FileName, getImgContentType(request.FileName))
	if err != nil {
		http.Error(w, "Error generating upload URL", http.StatusInternalServerError)
		return
	}

	imgUrl := strings.Split(uploadURL, "?")[0]
	err = s.Repository.UpdatePlayerProfileImage(userID, imgUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := struct {
		UploadURL string `json:"upload_url"`
		ImageURL  string `json:"image_url"`
	}{
		UploadURL: uploadURL,
		ImageURL:  strings.Split(uploadURL, "?")[0],
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}
func (s *Service) DeleteRanking(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rankingIDStr := r.PathValue("rankingId")
	rankingID, err := strconv.Atoi(rankingIDStr)
	if err != nil {
		http.Error(w, "Invalid ranking ID", http.StatusBadRequest)
		return
	}

	// Fetch ranking data before deletion to include in event
	ranking, gameID, otherPlayerIDs, err := s.Repository.GetRankingWithGamePlayers(uint(rankingID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Call the repository to delete the ranking
	err = s.Repository.DeleteRanking(uint(rankingID), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Publish ranking deleted event if the ranking had a player
	if ranking.PlayerID != nil {
		s.eventBus.Publish(events.RankingDeletedEvent{
			RankingID:      uint(rankingID),
			GameID:         gameID,
			PlayerID:       *ranking.PlayerID,
			OtherPlayerIDs: otherPlayerIDs,
			Date:           time.Now(),
		})
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) CreateDeck(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request CreateDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deck, err := s.Repository.CreateDeck(
		userID,
		request.Commander,
		request.Image,
		request.SecondaryImage,
		request.Crop,
		request.MoxfieldURL,
		request.Themes,
		request.Colors,
		request.Bracket,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to DTO (you may need to create a converter function)
	result := convertDeckToDto(deck)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

func (s *Service) GetPlayerDecks(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("playerId")
	if playerID == "" {
		http.Error(w, "Player ID is required", http.StatusBadRequest)
		return
	}

	p := pagination.ParsePagination(r)

	decks, total, err := s.Repository.GetPlayerDecks(playerID, p.PerPage, p.Offset())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to DTO
	items := make([]DeckResponse, 0, len(decks))
	for _, deck := range decks {
		items = append(items, convertDeckToDto(&deck))
	}

	result := pagination.PaginatedResult[DeckResponse]{
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

func (s *Service) GetPlayerGames(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("playerId")
	if playerID == "" {
		http.Error(w, "Player ID is required", http.StatusBadRequest)
		return
	}

	p := pagination.ParsePagination(r)

	games, total, err := s.Repository.GetPlayerGames(playerID, p.PerPage, p.Offset())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to DTO
	items := make([]GameResponse, 0, len(games))
	for _, game := range games {
		items = append(items, s.ConvertGameToDto(&game, true))
	}

	result := pagination.PaginatedResult[GameResponse]{
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

func (s *Service) UpdateMyPlayer(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var request UpdatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Build updates map
	updates := make(map[string]interface{})
	if request.MoxfieldUsername != nil {
		updates["moxfield_username"] = *request.MoxfieldUsername
	}

	if len(updates) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	player, err := s.Repository.UpdatePlayer(userID, updates)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := s.ConvertPlayerToResponse(player)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}
