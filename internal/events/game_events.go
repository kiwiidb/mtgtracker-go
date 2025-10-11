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
	GameID uint
	Date   time.Time
}

func (e GameDeletedEvent) EventName() string {
	return "game.deleted"
}
