package repository

import (
	"time"

	"gorm.io/gorm"
)

type Deck struct {
	Commander      string // in case of multi-name commanders, both names seperated by a /
	Image          string
	SecondaryImage string
	Crop           string
}

type Player struct {
	gorm.Model
	Name   string `gorm:"unique;not null"`
	Email  string `gorm:"unique;not null"`
	Image  string
	Groups []Group `gorm:"many2many:group_memberships;" json:"-"`
	Games  []Game  `gorm:"many2many:game_players;" json:"-"`
}

type Group struct {
	gorm.Model
	Name    string
	Image   string
	Players []Player `gorm:"many2many:group_memberships;"`
}

type Game struct {
	gorm.Model
	Duration *int
	Date     *time.Time
	Comments string
	Image    string
	Rankings []Ranking
	GroupID  uint
	Group    Group
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

	Player Player
	Deck   Deck `gorm:"embedded"`
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
