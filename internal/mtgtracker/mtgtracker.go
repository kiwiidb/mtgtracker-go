package mtgtracker

import (
	"encoding/json"
	"log"
	"mtgtracker/internal/repository"
	"mtgtracker/internal/scryfall"
	"net/http"
	"strconv"
)

type Service struct {
	Repository *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Repository: repo,
	}
}

func (s *Service) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /player/v1/signup", s.SignupPlayer)
	mux.HandleFunc("POST /player/v1/decks", s.AddDeckToPlayer)
	mux.HandleFunc("POST /player/v1/groups", s.CreateGroup)
	mux.HandleFunc("PUT /player/v1/groups/{groupId}/add/{email}", s.AddPlayerToGroup)
	mux.HandleFunc("PUT /player/v1/groups/{groupId}/games", s.AddGame)
	mux.HandleFunc("/player/v1/groups", s.GetGroups)
	mux.HandleFunc("/player/v1/games", s.GetGames)
	mux.HandleFunc("/player/v1/players/{playerId}/decks", s.GetDecks)
	mux.HandleFunc("DELETE /player/v1/decks/{deckId}", s.DeleteDeck)
	mux.HandleFunc("/player/v1/groups/{groupId}/ranking", s.GetRanking)
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

func (s *Service) AddDeckToPlayer(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var request AddDeckToPlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the repository to add the deck to the player
	imgUris, err := scryfall.GetCardImageURIs(request.Commander)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	crop := imgUris.ArtCrop
	img := imgUris.Normal

	deck, err := s.Repository.AddDeckToPlayer(request.PlayerID, request.MoxfieldURL, request.Commander, img, crop)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the created deck as a JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(deck)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

func (s *Service) CreateGroup(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var request CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the repository to create the group
	group, err := s.Repository.CreateGroup(request.CreatorID, request.Name, request.Image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the created group as a JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(group)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}
func (s *Service) AddPlayerToGroup(w http.ResponseWriter, r *http.Request) {
	// Call the repository to add the player to the group
	groupID := r.PathValue("groupId")
	groupIDInt, err := strconv.Atoi(groupID)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}
	email := r.PathValue("email")
	err = s.Repository.AddPlayerToGroup(uint(groupIDInt), email)
	if err != nil {
		log.Println("Error adding player to group:", err, email)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
		rankings = append(rankings, repository.Ranking{
			PlayerID:     rank.PlayerID,
			DeckID:       rank.DeckID,
			Position:     rank.Position,
			CouldHaveWon: rank.CouldHaveWon,
		})
	}
	game, err := s.Repository.InsertGame(request.GroupID, request.Duration, request.Comments, request.Image, request.Date, rankings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(game)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}

func (s *Service) GetGroups(w http.ResponseWriter, r *http.Request) {
	// Call the repository to get the groups
	groups, err := s.Repository.GetGroups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the groups as a JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(groups)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
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
func (s *Service) GetDecks(w http.ResponseWriter, r *http.Request) {
	// Call the repository to get the decks
	playerID := r.PathValue("playerId")
	playerIDInt, err := strconv.Atoi(playerID)
	if err != nil {
		http.Error(w, "Invalid player ID", http.StatusBadRequest)
		return
	}
	decks, err := s.Repository.GetDecks(uint(playerIDInt))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the decks as a JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(decks)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}
func (s *Service) GetRanking(w http.ResponseWriter, r *http.Request) {
	// Call the repository to get the ranking
	groupID := r.PathValue("groupId")
	groupIDInt, err := strconv.Atoi(groupID)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}
	ranking, err := s.Repository.GetGroupRankingByWins(uint(groupIDInt))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the ranking as a JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ranking)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}
