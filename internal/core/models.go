package core

import (
	"time"

	"gorm.io/gorm"
)

const (
	EventTypeInit      = "init"
	EventTypeIncrement = "increment"
	EventTypeDecrement = "decrement"
	EventTypeImage     = "image"
	EventTypeScoop     = "scoop"
)

type Deck struct {
	gorm.Model
	MoxfieldURL    *string  `json:"moxfield_url"`
	Themes         []string `gorm:"serializer:json" json:"themes"`
	Bracket        *uint    `json:"bracket"`
	Commander      string   `json:"commander"`
	Colors         []string `gorm:"serializer:json" json:"colors"` // Scryfall color codes: W, U, B, R, G, C
	Image          string   `json:"image"`
	SecondaryImage string   `json:"secondary_image"`
	Crop           string   `json:"crop"`
	PlayerID       *string  `json:"player_id,omitempty"`
	GameCount      int      `gorm:"default:0" json:"game_count"`
	WinCount       int      `gorm:"default:0" json:"win_count"`

	Player *Player `gorm:"foreignKey:PlayerID;references:FirebaseID" json:"player,omitempty"`
}

type SimpleDeck struct {
	Commander      string `json:"commander"`
	Image          string `json:"image"`
	SecondaryImage string `json:"secondary_image"`
	Crop           string `json:"crop"`
}

type Player struct {
	FirebaseID       string `gorm:"primaryKey" json:"firebase_id"`
	Name             string `gorm:"unique;not null" json:"name"`
	Email            string `gorm:"unique;not null" json:"email"`
	Image            string
	MoxfieldUsername string `json:"moxfield_username"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	Games            []Game         `gorm:"-" json:"-"` // Not a GORM relationship, populated manually
	Decks            []Deck         `gorm:"foreignKey:PlayerID;references:FirebaseID" json:"decks,omitempty"`
}
type Game struct {
	gorm.Model
	CreatorID  *string `json:"creator_id,omitempty"`
	Duration   *int
	Date       *time.Time
	Comments   string
	Image      string
	Rankings   []Ranking
	Finished   bool
	GameEvents []GameEvent // Add relation: a game has many game events

	Creator *Player `gorm:"foreignKey:CreatorID;references:FirebaseID" json:"creator,omitempty"`
}

type Ranking struct {
	gorm.Model
	GameID         uint    `json:"game_id"`
	PlayerID       *string `json:"player_id,omitempty"`
	DeckID         *uint   `json:"deck_id,omitempty"` // Reference to player's deck (optional)
	Position       int     `json:"position"`
	CouldHaveWon   bool    `json:"could_have_won"`
	EarlySolRing   bool    `json:"early_sol_ring"`
	StartingPlayer bool    `json:"starting_player"`
	PlayerName     string  `gorm:"-"`

	Player       *Player    `gorm:"foreignKey:PlayerID;references:FirebaseID" json:"player,omitempty"`
	Deck         *Deck      `gorm:"foreignKey:DeckID;references:ID" json:"deck,omitempty"` // Reference to Deck model
	DeckEmbedded SimpleDeck `gorm:"embedded" json:"deck_embedded,omitempty"`               // Embedded deck info for games without deck reference
}

type DeckWin struct {
	DeckID uint
	Deck   Deck
	Wins   int
}

type PlayerWin struct {
	PlayerID string
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
	SourceRankingID      *uint   // Made nullable with pointer
	TargetRankingID      *uint   // Made nullable with pointer
	ImageUrl             string  // New field for uploaded image URL
	Comment              *string // New field for text description

	SourceRanking *Ranking `gorm:"foreignKey:SourceRankingID;references:ID"` // Made nullable with pointer
	TargetRanking *Ranking `gorm:"foreignKey:TargetRankingID;references:ID"` // Made nullable with pointer
}
