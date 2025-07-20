package mtgtracker

import (
	"encoding/json"
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
	mux.HandleFunc("PUT /game/v1/games/{gameId}", s.UpdateGame)           // new
	mux.HandleFunc("GET /game/v1/games/{gameId}", s.GetGame)              // new
	mux.HandleFunc("DELETE /game/v1/games/{gameId}", s.DeleteGame)        // new
	mux.HandleFunc("POST /game/v1/games/{gameId}/events", s.AddGameEvent) // new
	mux.HandleFunc("GET /ranking/v1/rankings/pending", s.GetPendingRankings)
	mux.HandleFunc("PUT /ranking/v1/rankings/{rankingId}/accept", s.AcceptRanking)
	mux.HandleFunc("PUT /ranking/v1/rankings/{rankingId}/decline", s.DeclineRanking)
}

func (s *Service) GetPendingRankings(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the context
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Call the repository to get the rankings to accept
	rankings, err := s.Repository.GetPendingRankings(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(rankings)
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
	err = s.Repository.AcceptRanking(uint(rankingIDInt))
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
	email := ""
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
	for i, player := range players {
		// Convert the player to a DTO
		result[i] = convertPlayerToDto(&player)
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

func (s *Service) GetPlayer(w http.ResponseWriter, r *http.Request) {
	playerID := r.PathValue("playerId")
	playerIDInt, err := strconv.Atoi(playerID)
	if err != nil {
		http.Error(w, "Invalid player ID", http.StatusBadRequest)
		return
	}
	// Call the repository to get the player
	player, err := s.Repository.GetPlayerByID(uint(playerIDInt))
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

func (s *Service) DeleteDeck(w http.ResponseWriter, r *http.Request) {
	// Call the repository to delete the deck
	deckID := r.PathValue("deckId")
	deckIDInt, err := strconv.Atoi(deckID)
	if err != nil {
		http.Error(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}
	err = s.Repository.DeleteDeck(uint(deckIDInt))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) CreateGame(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
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

		rankings = append(rankings, repository.Ranking{
			PlayerID: rank.PlayerID,
			Position: rank.Position,
			Deck:     convertDeck(rank.Deck),
		})
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
	for i, game := range games {
		// Convert the game to a DTO
		result[i] = convertGameToDto(&game)
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
	newRankings := []repository.Ranking{}
	for _, rank := range request.Rankings {
		// the only thing that can change is the position
		for _, oldRank := range rankings {
			if oldRank.PlayerID == rank.PlayerID {
				oldRank.Position = rank.Position
				newRankings = append(newRankings, oldRank)
				break
			}
		}
	}
	if len(newRankings) != len(rankings) {
		http.Error(w, "Invalid rankings", http.StatusBadRequest)
		return
	}

	// Call the repository to update the game (implement UpdateGame in your repository)
	updatedGame, err := s.Repository.UpdateGame(uint(gameId), newRankings, request.Finished)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedGame.GameEvents = []repository.GameEvent{} // Clear event
	result := convertGameToDto(updatedGame)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}
