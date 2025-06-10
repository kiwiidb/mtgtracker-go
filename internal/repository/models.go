package repository

import (
	"time"

	"gorm.io/gorm"
)

type Deck struct {
	Commander      string `json:"commander"`
	Image          string `json:"image"`
	SecondaryImage string `json:"secondary_image"`
	Crop           string `json:"crop"`
}

type Player struct {
	gorm.Model
	Name  string `gorm:"unique;not null"`
	Email string `gorm:"unique;not null"`
	Image string
	Games []Game `gorm:"many2many:game_players;" json:"-"`
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
	GameID         uint
	PlayerID       uint
	Position       int
	CouldHaveWon   bool
	EarlySolRing   bool
	StartingPlayer bool
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
	SourceRankingID      uint
	TargetRankingID      uint

	SourceRanking Ranking `gorm:"foreignKey:SourceRankingID;references:ID"`
	TargetRanking Ranking `gorm:"foreignKey:TargetRankingID;references:ID"`
}
