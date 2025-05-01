package repository

import (
	"time"

	"gorm.io/gorm"
)

// Models

type Deck struct {
	gorm.Model
	MoxfieldURL string
	Commander   string
	Image       string
	PlayerID    uint
}

type Player struct {
	gorm.Model
	Name   string `gorm:"unique;not null"`
	Email  string `gorm:"unique;not null"`
	Image  string
	Decks  []Deck
	Groups []Group `gorm:"many2many:group_memberships;"`
	Games  []Game  `gorm:"many2many:game_players;"`
}

type Group struct {
	gorm.Model
	Name    string
	Image   string
	Players []Player `gorm:"many2many:group_memberships;"`
}

type Game struct {
	gorm.Model
	Duration int
	Date     *time.Time
	Comments string
	Image    string
	Rankings []Ranking
	GroupID  uint
	Group    Group
}

type Ranking struct {
	gorm.Model
	GameID   uint
	PlayerID uint
	DeckID   uint
	Position int

	Player Player
	Deck   Deck
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
