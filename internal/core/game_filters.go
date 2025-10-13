package core

import "gorm.io/gorm"

// GameFilter represents search criteria for games
type GameFilter struct {
	PlayerIDs  []string // Games where ANY of these players participated
	Commanders []string // Games where ANY of these commanders were played
	AllPlayers []string // Games where ALL of these players participated (AND)
}

// ApplyGameFilters applies the filter criteria to a GORM query
func ApplyGameFilters(db *gorm.DB, filter GameFilter) *gorm.DB {
	query := db

	// Filter by player IDs (OR condition - any of these players)
	if len(filter.PlayerIDs) > 0 {
		subQuery := db.Session(&gorm.Session{}).Model(&Ranking{}).
			Select("game_id").
			Where("player_id IN ?", filter.PlayerIDs)
		query = query.Where("id IN (?)", subQuery)
	}

	// Filter by commanders (OR condition - any of these commanders)
	if len(filter.Commanders) > 0 {
		// Search in both DeckEmbedded.Commander (JSONB) and referenced Deck.Commander
		subQuery := db.Session(&gorm.Session{}).Model(&Ranking{}).
			Select("DISTINCT game_id").
			Joins("LEFT JOIN decks ON rankings.deck_id = decks.id").
			Where("rankings.deck_embedded->>'commander' IN ? OR decks.commander IN ?",
				filter.Commanders, filter.Commanders)
		query = query.Where("id IN (?)", subQuery)
	}

	// Filter by all players (AND condition - all of these players must be in the game)
	if len(filter.AllPlayers) > 0 {
		// For each player, ensure they're in the game
		for _, playerID := range filter.AllPlayers {
			subQuery := db.Session(&gorm.Session{}).Model(&Ranking{}).
				Select("game_id").
				Where("player_id = ?", playerID)
			query = query.Where("id IN (?)", subQuery)
		}
	}

	return query
}

// SearchGamesWithFilters searches games with complex filtering
func (r *Repository) SearchGamesWithFilters(filter GameFilter, limit, offset int) ([]Game, int64, error) {
	var games []Game
	var total int64

	// Build the base query
	baseQuery := r.DB.Model(&Game{})
	baseQuery = ApplyGameFilters(baseQuery, filter)

	// Get total count
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with preloads
	query := r.DB.Model(&Game{})
	query = ApplyGameFilters(query, filter)

	err := query.
		Preload("Rankings", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Player").Preload("Deck")
		}).
		Preload("GameEvents").
		Preload("Creator").
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&games).Error

	if err != nil {
		return nil, 0, err
	}

	return games, total, nil
}
