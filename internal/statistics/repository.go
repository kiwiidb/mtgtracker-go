package statistics

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	err := db.AutoMigrate(&PlayerStats{})
	if err != nil {
		log.Fatalf("Failed to migrate stats repo: %v", err)
	}
	return &Repository{DB: db}
}

// GetPlayerStatsTimeSeries retrieves all statistics for a specific player ordered by timestamp
func (r *Repository) GetPlayerStatsTimeSeries(playerID string, limit, offset int) ([]PlayerStats, int64, error) {
	var stats []PlayerStats
	var total int64

	// Get total count for this player
	if err := r.DB.Model(&PlayerStats{}).Where("player_id = ?", playerID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results ordered by timestamp descending (most recent first)
	err := r.DB.Where("player_id = ?", playerID).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&stats).Error
	if err != nil {
		return nil, 0, err
	}

	return stats, total, nil
}

// GetLatestPlayerStats retrieves the most recent statistics for a specific player
func (r *Repository) GetLatestPlayerStats(playerID string) (*PlayerStats, error) {
	var stats PlayerStats
	err := r.DB.Where("player_id = ?", playerID).
		Order("timestamp DESC").
		First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// CreatePlayerStats creates a new statistics entry
func (r *Repository) CreatePlayerStats(stats *PlayerStats) error {
	if stats.Timestamp.IsZero() {
		stats.Timestamp = time.Now()
	}
	return r.DB.Create(stats).Error
}

// GetAllLatestPlayerStats retrieves the most recent statistics for all players with pagination
func (r *Repository) GetAllLatestPlayerStats(limit, offset int) ([]PlayerStats, int64, error) {
	var stats []PlayerStats

	// Subquery to get the latest timestamp for each player
	subQuery := r.DB.Model(&PlayerStats{}).
		Select("player_id, MAX(timestamp) as max_timestamp").
		Group("player_id")

	// Get total count of unique players
	var total int64
	if err := r.DB.Table("(?) as sub", subQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Join with the subquery to get full records
	err := r.DB.Table("player_stats").
		Select("player_stats.*").
		Joins("INNER JOIN (?) as latest ON player_stats.player_id = latest.player_id AND player_stats.timestamp = latest.max_timestamp", subQuery).
		Order("player_stats.elo DESC").
		Limit(limit).
		Offset(offset).
		Scan(&stats).Error
	if err != nil {
		return nil, 0, err
	}

	return stats, total, nil
}

// DeletePlayerStats removes all statistics entries for a specific player
func (r *Repository) DeletePlayerStats(playerID string) error {
	return r.DB.Where("player_id = ?", playerID).Delete(&PlayerStats{}).Error
}
