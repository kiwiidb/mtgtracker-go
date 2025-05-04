package repository

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) GetDecks(playerId uint) ([]Deck, error) {
	var decks []Deck
	err := r.DB.Where(Deck{PlayerID: playerId}).Find(&decks).Error
	if err != nil {
		return nil, err
	}
	return decks, nil
}

func (r *Repository) GetGames() ([]Game, error) {
	var games []Game
	err := r.DB.Preload("Rankings").Find(&games).Error
	if err != nil {
		return nil, err
	}
	return games, nil
}

func (r *Repository) DeleteDeck(deckID uint) error {
	var deck Deck
	if err := r.DB.First(&deck, deckID).Error; err != nil {
		return errors.New("deck not found")
	}
	return r.DB.Delete(&deck).Error
}

func (r *Repository) GetGroups() ([]Group, error) {
	var groups []Group
	// preload the players for the groups
	// preload the decks for the players
	err := r.DB.Preload("Players").Preload("Players.Decks").Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func NewRepository(db *gorm.DB) *Repository {
	err := db.AutoMigrate(&Player{}, &Deck{}, &Group{}, &Game{}, &Ranking{})
	if err != nil {
		log.Fatal(err)
	}
	return &Repository{
		DB: db,
	}
}
func (r *Repository) InsertPlayer(name string, email string, image string) (*Player, error) {
	player := Player{Name: name, Email: email, Image: image}
	result := r.DB.Create(&player)
	return &player, result.Error
}

func (r *Repository) AddDeckToPlayer(playerID uint, moxfieldURL, commander, img, secondaryImg string) (*Deck, error) {
	deck := Deck{
		MoxfieldURL:    moxfieldURL,
		Commander:      commander,
		Image:          img,
		PlayerID:       playerID,
		SecondaryImage: secondaryImg,
	}
	result := r.DB.Create(&deck)
	return &deck, result.Error
}

func (r *Repository) CreateGroup(creatorID uint, name string, image string) (*Group, error) {
	var creator Player
	if err := r.DB.First(&creator, creatorID).Error; err != nil {
		return nil, err
	}

	group := Group{Name: name, Image: image, Players: []Player{creator}}
	result := r.DB.Create(&group)
	return &group, result.Error
}

func (r *Repository) AddPlayerToGroup(groupID uint, email string) error {
	//find the player by email
	var player Player
	err := r.DB.Where("email = ?", email).First(&player).Error
	if err != nil {
		return errors.New("player not found")
	}
	var group Group
	if err := r.DB.Preload("Players").First(&group, groupID).Error; err != nil {
		return err
	}
	return r.DB.Model(&group).Association("Players").Append(&player)
}

func (r *Repository) InsertGame(groupID uint, duration int, comments, image string, date *time.Time, rankings []Ranking) (*Game, error) {
	var group Group
	if err := r.DB.First(&group, groupID).Error; err != nil {
		return nil, errors.New("invalid group ID")
	}
	// Ensure each ranking has valid player and deck (optional but safe)
	for _, rank := range rankings {
		var player Player
		if err := r.DB.First(&player, rank.PlayerID).Error; err != nil {
			return nil, errors.New("invalid player ID")
		}
		var deck Deck
		if err := r.DB.First(&deck, rank.DeckID).Error; err != nil {
			return nil, errors.New("invalid deck ID")
		}
	}

	game := Game{
		GroupID:  groupID,
		Duration: duration,
		Date:     date,
		Comments: comments,
		Image:    image,
		Rankings: rankings,
	}

	if err := r.DB.Create(&game).Error; err != nil {
		return nil, err
	}

	return &game, nil
}
