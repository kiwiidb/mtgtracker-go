package statistics

import (
	"time"

	"gorm.io/gorm"
)

type PlayerStats struct {
	gorm.Model
	PlayerID        string    `gorm:"index;not null" json:"player_id"`
	Timestamp       time.Time `gorm:"index;not null" json:"timestamp"`
	TotalWins       int       `json:"total_wins"`
	Winrate         float64   `json:"winrate"`
	RollingWinrate  float64   `json:"rolling_winrate"`   // Winrate over last N games (moving average)
	GameCount       int       `json:"game_count"`
	GameDuration    int       `json:"game_duration"`     // Average game duration in minutes
	Streak          int       `json:"streak"`            // Current win/loss streak (positive = wins, negative = losses)
	Elo             int       `json:"elo"`
}
