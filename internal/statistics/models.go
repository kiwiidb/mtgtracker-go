package statistics

import (
	"time"

	"gorm.io/gorm"
)

type PlayerStats struct {
	gorm.Model
	PlayerID       string    `gorm:"index;not null"`
	Timestamp      time.Time `gorm:"index;not null"`
	TotalWins      int
	Winrate        float64
	RollingWinrate float64 // Winrate over last N games (moving average)
	GameCount      int
	GameDuration   int // Average game duration in minutes
	Streak         int // Current win/loss streak (positive = wins, negative = losses)
	Elo            int
}

// PlayerStatsResponse is the DTO for API responses with snake_case JSON tags
type PlayerStatsResponse struct {
	ID             uint      `json:"id"`
	PlayerID       string    `json:"player_id"`
	Timestamp      time.Time `json:"timestamp"`
	TotalWins      int       `json:"total_wins"`
	Winrate        float64   `json:"winrate"`
	RollingWinrate float64   `json:"rolling_winrate"`
	GameCount      int       `json:"game_count"`
	GameDuration   int       `json:"game_duration"`
	Streak         int       `json:"streak"`
	Elo            int       `json:"elo"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ToResponse converts PlayerStats to PlayerStatsResponse
func (ps *PlayerStats) ToResponse() PlayerStatsResponse {
	return PlayerStatsResponse{
		ID:             ps.ID,
		PlayerID:       ps.PlayerID,
		Timestamp:      ps.Timestamp,
		TotalWins:      ps.TotalWins,
		Winrate:        ps.Winrate,
		RollingWinrate: ps.RollingWinrate,
		GameCount:      ps.GameCount,
		GameDuration:   ps.GameDuration,
		Streak:         ps.Streak,
		Elo:            ps.Elo,
		CreatedAt:      ps.CreatedAt,
		UpdatedAt:      ps.UpdatedAt,
	}
}
