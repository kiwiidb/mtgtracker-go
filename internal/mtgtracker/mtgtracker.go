package mtgtracker

import (
	"encoding/json"
	"fmt"
	"log"
	"mtgtracker/internal/repository"
	"mtgtracker/internal/scryfall"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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
	mux.HandleFunc("POST /player/v1/groups", s.CreateGroup)
	mux.HandleFunc("PUT /player/v1/groups/{groupId}/add/{email}", s.AddPlayerToGroup)
	mux.HandleFunc("POST /game/v1/games", s.AddGame)
	mux.HandleFunc("/game/v1/games", s.GetGames)
	mux.HandleFunc("/player/v1/groups", s.GetGroups)
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
