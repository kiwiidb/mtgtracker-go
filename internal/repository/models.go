package repository

import (
	"time"

	"gorm.io/gorm"
)

const (
	EventTypeIncrement = "increment"
	EventTypeDecrement = "decrement"
	EventTypeImage     = "image"
	EventTypeScoop     = "scoop"
)

type Deck struct {
	Commander      string `json:"commander"`
	Image          string `json:"image"`
	SecondaryImage string `json:"secondary_image"`
	Crop           string `json:"crop"`
}

type Player struct {
	gorm.Model
	FirebaseID string `gorm:"unique;not null" json:"firebase_id"`
	Name       string `gorm:"unique;not null" json:"name"`
	Email      string `gorm:"unique;not null" json:"email"`
	Image      string
	Games      []Game `gorm:"-" json:"-"` // Not a GORM relationship, populated manually
}
type Game struct {
	gorm.Model
	Duration   *int
	Date       *time.Time
	Comments   string
	Image      string
	Rankings   []Ranking
	Finished   bool
	GameEvents []GameEvent // Add relation: a game has many game events
}

type Ranking struct {
	gorm.Model
	GameID         uint   `json:"game_id"`
	PlayerID       uint   `json:"player_id"`
	Position       int    `json:"position"`
	CouldHaveWon   bool   `json:"could_have_won"`
	EarlySolRing   bool   `json:"early_sol_ring"`
	StartingPlayer bool   `json:"starting_player"`
	PlayerName     string `gorm:"-"`

	Player Player `json:"player"`
	Deck   Deck   `gorm:"embedded" json:"deck"` // Use embedded struct for Deck
}

type DeckWin struct {
	DeckID uint
	Deck   Deck
	Wins   int
}

type PlayerWin struct {
	PlayerID uint
	Player   Player
	Wins     int
	DeckWins []DeckWin
}

type GameEvent struct {
	gorm.Model
	GameID               uint
	EventType            string
	DamageDelta          int
	TargetLifeTotalAfter int
	SourceRankingID      *uint  // Made nullable with pointer
	TargetRankingID      *uint  // Made nullable with pointer
	ImageUrl             string // New field for uploaded image URL

	SourceRanking *Ranking `gorm:"foreignKey:SourceRankingID;references:ID"` // Made nullable with pointer
	TargetRanking *Ranking `gorm:"foreignKey:TargetRankingID;references:ID"` // Made nullable with pointer
}
