package events

import "time"

// GameCreatedEvent is published when a new game is created
type GameCreatedEvent struct {
	GameID     uint
	CreatorID  string
	RankingIDs []uint
	Date       time.Time
}

func (e GameCreatedEvent) EventName() string {
	return "game.created"
}

// GameFinishedEvent is published when a game is marked as finished
type GameFinishedEvent struct {
	GameID     uint
	RankingIDs []uint
	Date       time.Time
}

func (e GameFinishedEvent) EventName() string {
	return "game.finished"
}

// GameDeletedEvent is published when a game is deleted
type GameDeletedEvent struct {
	GameID     uint
	RankingIDs []uint
	PlayerIDs  []string // Player IDs from rankings (for follow count decrements)
	Date       time.Time
}

func (e GameDeletedEvent) EventName() string {
	return "game.deleted"
}

// RankingDeletedEvent is published when a player removes themselves from a game
type RankingDeletedEvent struct {
	RankingID uint
	GameID    uint
	PlayerID  string   // The player who was removed
	OtherPlayerIDs []string // Other players in the game (for follow count decrements)
	Date      time.Time
}

func (e RankingDeletedEvent) EventName() string {
	return "ranking.deleted"
}
