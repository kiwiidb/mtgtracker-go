package mtgtracker

import (
	"encoding/json"
	"fmt"
	"log"
	"mtgtracker/internal/repository"
	"mtgtracker/internal/scryfall"
	"net/http"
	"path/filepath"
	"regexp"
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
	mux.HandleFunc("POST /game/v1/games", s.AddGame)
	mux.HandleFunc("GET /game/v1/games", s.GetGames)
	mux.HandleFunc("POST /game/v1/games/{gameId}/events", s.AddGameEvent) // new
	mux.HandleFunc("PUT /game/v1/games/{gameId}", s.UpdateGame)           // new
	mux.HandleFunc("GET /game/v1/games/{gameId}", s.GetGame)              // new
	mux.HandleFunc("DELETE /game/v1/games/{gameId}", s.DeleteGame)        // new
}

func (s *Service) SignupPlayer(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var request SignupPlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the repository to insert the player
	player, err := s.Repository.InsertPlayer(request.Name, request.Email, request.Image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the created player as a JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(player)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

func (s *Service) GetPlayers(w http.ResponseWriter, r *http.Request) {
	// Call the repository to get the players
	players, err := s.Repository.GetPlayers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the players as a JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(players)
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
	// Return the player as a JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Error encoding response:", err)
		return
	}
}

func findCommander(name string) (*repository.Deck, error) {
	// Call the repository to add the deck to the player
	commanderCard, err := scryfall.GetCard(name)
	if err != nil {
		return nil, err
	}
	commander := commanderCard.Name
	img := commanderCard.ImageURIs.Normal
	crop := commanderCard.ImageURIs.ArtCrop
	secondaryImg := ""
	// first check if the commander is a double-faced card
	if len(commanderCard.CardFaces) > 1 {
		img = commanderCard.CardFaces[0].ImageURIs.Normal
		crop = commanderCard.CardFaces[0].ImageURIs.ArtCrop
		secondaryImg = commanderCard.CardFaces[1].ImageURIs.Normal
		commander = commanderCard.CardFaces[0].Name + "/" + commanderCard.CardFaces[1].Name
	}
	// then check if the commander has a partner
	if partner, ok := findPartner(commanderCard.OracleText); ok {
		commander = strings.Join([]string{commander, partner.Name}, "/")
		secondaryImg = partner.ImageURIs.Normal
	}
	return &repository.Deck{
		Commander:      commander,
		Image:          img,
		Crop:           crop,
		SecondaryImage: secondaryImg,
	}, nil
}

func findPartner(oracleText string) (*scryfall.Card, bool) {
	// Use a regex to find "Partner with <partner name>"
	re := regexp.MustCompile(`Partner with ([^\n]+)`)
	matches := re.FindStringSubmatch(oracleText)

	if len(matches) < 2 {
		return nil, false // No partner found
	}

	partnerName := matches[1]
	log.Println("Partner found:", partnerName)
	partnerCard, err := scryfall.GetCard(partnerName)
	if err != nil {
		log.Printf("Failed to fetch partner card: %v", err)
		return nil, false
	}

	return partnerCard, true
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

func (s *Service) AddGame(w http.ResponseWriter, r *http.Request) {
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
		// first, find the full name of the commander
		deck, err := findCommander(rank.Commander)
		if err != nil {
			fmt.Println("Error finding commander:", err, rank.Commander)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rankings = append(rankings, repository.Ranking{
			PlayerID:       rank.PlayerID,
			Position:       rank.Position,
			CouldHaveWon:   rank.CouldHaveWon,
			EarlySolRing:   rank.EarlySolRing,
			StartingPlayer: rank.StartingPlayer,
			Deck:           *deck,
		})
	}
	game, err := s.Repository.InsertGame(request.Duration, request.Comments, request.Image, request.Date, request.Finished, rankings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := convertGameToDto(game)
	w.Header().Set("Content-Type", "application/json")
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

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(eventDto)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

func getImgContentType(s string) string {
	switch filepath.Ext(s) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
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
	w.Header().Set("Content-Type", "application/json")
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

	// Return the games as a JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(games)
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

	w.Header().Set("Content-Type", "application/json")
	updatedGame.GameEvents = []repository.GameEvent{} // Clear event
	err = json.NewEncoder(w).Encode(updatedGame)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}
