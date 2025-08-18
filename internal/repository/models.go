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
	FirebaseID string `gorm:"primaryKey" json:"firebase_id"`
	Name       string `gorm:"unique;not null" json:"name"`
	Email      string `gorm:"unique;not null" json:"email"`
	Image      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
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

type RankingStatus string

const (
	StatusPending  RankingStatus = "pending"
	StatusAccepted RankingStatus = "accepted"
	StatusDeclined RankingStatus = "declined"
)

type Ranking struct {
	gorm.Model
	GameID         uint          `json:"game_id"`
	PlayerID       *string       `json:"player_id,omitempty"`
	Position       int           `json:"position"`
	CouldHaveWon   bool          `json:"could_have_won"`
	EarlySolRing   bool          `json:"early_sol_ring"`
	StartingPlayer bool          `json:"starting_player"`
	Status         RankingStatus `json:"status"`
	PlayerName     string        `gorm:"-"`

	Player *Player `gorm:"foreignKey:PlayerID;references:FirebaseID" json:"player,omitempty"`
	Deck   Deck   `gorm:"embedded" json:"deck"` // Use embedded struct for Deck
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
	SourceRankingID      *uint  // Made nullable with pointer
	TargetRankingID      *uint  // Made nullable with pointer
	ImageUrl             string // New field for uploaded image URL

	SourceRanking *Ranking `gorm:"foreignKey:SourceRankingID;references:ID"` // Made nullable with pointer
	TargetRanking *Ranking `gorm:"foreignKey:TargetRankingID;references:ID"` // Made nullable with pointer
}

type Follow struct {
	gorm.Model
	Player1ID string `gorm:"not null" json:"player1_id"`
	Player2ID string `gorm:"not null" json:"player2_id"`

	Player1 Player `gorm:"foreignKey:Player1ID;references:FirebaseID" json:"player1"`
	Player2 Player `gorm:"foreignKey:Player2ID;references:FirebaseID" json:"player2"`
}
