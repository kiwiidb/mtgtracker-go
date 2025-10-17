package statistics

import (
	"log"
	"math"
	"mtgtracker/internal/core"
	"mtgtracker/internal/events"
)

type CoreService interface {
	GetGameByID(gameID uint) (*core.Game, error)
}

// EventHandlers manages event subscriptions for the statistics package
type EventHandlers struct {
	repo        *Repository
	coreService CoreService
}

// NewEventHandlers creates a new event handler instance
func NewEventHandlers(repo *Repository, coreService CoreService) *EventHandlers {
	return &EventHandlers{
		repo:        repo,
		coreService: coreService,
	}
}

// RegisterHandlers subscribes to all relevant events
func (h *EventHandlers) RegisterHandlers(bus *events.EventBus) {
	bus.Subscribe("game.finished", h.HandleGameFinished)
	log.Println("Statistics event handlers registered")
}

// HandleGameFinished processes game finished events and updates player statistics
func (h *EventHandlers) HandleGameFinished(event events.Event) error {
	e, ok := event.(events.GameFinishedEvent)
	if !ok {
		log.Printf("Invalid event type for game.finished: %T", event)
		return nil
	}

	log.Printf("Processing game.finished event for statistics (game %d)", e.GameID)

	// Fetch game data
	game, err := h.coreService.GetGameByID(e.GameID)
	if err != nil {
		log.Printf("Failed to fetch game %d: %v", e.GameID, err)
		return err
	}

	// Update statistics for each player
	for _, ranking := range game.Rankings {
		if ranking.PlayerID == nil {
			continue // Skip rankings without players
		}

		playerID := *ranking.PlayerID

		// Get current latest stats for the player
		currentStats, err := h.repo.GetLatestPlayerStats(playerID)
		if err != nil {
			// If no stats exist, create initial stats
			currentStats = &PlayerStats{
				PlayerID:       playerID,
				TotalWins:      0,
				Winrate:        0,
				RollingWinrate: 0,
				GameCount:      0,
				GameDuration:   0,
				Streak:         0,
				Elo:            1000, // Starting ELO
			}
		}

		// Get all player stats for ELO calculation
		allPlayerStats := make(map[string]*PlayerStats)
		allPlayerStats[playerID] = currentStats
		for _, r := range game.Rankings {
			if r.PlayerID != nil && *r.PlayerID != playerID {
				otherStats, err := h.repo.GetLatestPlayerStats(*r.PlayerID)
				if err != nil {
					// Initialize with default ELO if player has no stats
					otherStats = &PlayerStats{
						PlayerID: *r.PlayerID,
						Elo:      1000,
					}
				}
				allPlayerStats[*r.PlayerID] = otherStats
			}
		}

		// Calculate new stats
		newStats := h.calculateNewStats(currentStats, &ranking, game, allPlayerStats)

		// Create new stats entry
		err = h.repo.CreatePlayerStats(newStats)
		if err != nil {
			log.Printf("Failed to create stats for player %s: %v", playerID, err)
			// Continue processing other players
		} else {
			log.Printf("Updated stats for player %s: ELO %d, Winrate %.2f%%",
				playerID, newStats.Elo, newStats.Winrate*100)
		}
	}

	return nil
}

// calculateNewStats computes updated statistics based on game result
func (h *EventHandlers) calculateNewStats(current *PlayerStats, ranking *core.Ranking, game *core.Game, allPlayerStats map[string]*PlayerStats) *PlayerStats {
	won := ranking.Position == 1
	newGameCount := current.GameCount + 1

	// Update total wins
	newTotalWins := current.TotalWins
	if won {
		newTotalWins++
	}

	// Calculate new overall winrate
	newWinrate := float64(newTotalWins) / float64(newGameCount)

	// Calculate rolling winrate (moving average over last 10 games)
	newRollingWinrate := h.calculateRollingWinrate(current.PlayerID, won)

	// Update streak
	newStreak := h.calculateStreak(current.Streak, won)

	// Update total game duration
	newGameDuration := current.GameDuration
	if game.Duration != nil {
		newGameDuration = current.GameDuration + *game.Duration
	}

	// Calculate new ELO based on performance against all other players
	newElo := h.calculateMultiplayerElo(current.PlayerID, ranking, game, allPlayerStats)

	// Ensure streak doesn't go negative
	if newStreak < 0 {
		newStreak = 0
	}

	return &PlayerStats{
		PlayerID:       current.PlayerID,
		TotalWins:      newTotalWins,
		Winrate:        newWinrate,
		RollingWinrate: newRollingWinrate,
		GameCount:      newGameCount,
		GameDuration:   newGameDuration,
		Streak:         newStreak,
		Elo:            newElo,
	}
}

// calculateMultiplayerElo computes ELO rating by comparing against each opponent
// Uses the formula: R'i = Ri + K/(N-1) * Σ(Sij - Eij) for all j ≠ i
func (h *EventHandlers) calculateMultiplayerElo(playerID string, ranking *core.Ranking, game *core.Game, allPlayerStats map[string]*PlayerStats) int {
	currentElo := allPlayerStats[playerID].Elo
	kFactor := 32.0
	numPlayers := len(game.Rankings)

	if numPlayers <= 1 {
		return currentElo // No change if playing alone
	}

	// Create a map of playerID to their ranking position
	playerPositions := make(map[string]int)
	for _, r := range game.Rankings {
		if r.PlayerID != nil {
			playerPositions[*r.PlayerID] = r.Position
		}
	}

	currentPosition := ranking.Position
	totalChange := 0.0

	// Compare against each other player
	for _, otherRanking := range game.Rankings {
		if otherRanking.PlayerID == nil || *otherRanking.PlayerID == playerID {
			continue // Skip self
		}

		otherPlayerID := *otherRanking.PlayerID
		otherElo := allPlayerStats[otherPlayerID].Elo
		otherPosition := otherRanking.Position

		// Sij: 1 if player i ranked higher (lower position number) than j, else 0
		var sij float64
		if currentPosition < otherPosition {
			sij = 1.0
		} else {
			sij = 0.0
		}

		// Eij: Expected probability of player i beating player j
		// Eij = 1 / (1 + 10^((Rj - Ri)/400))
		eij := 1.0 / (1.0 + math.Pow(10.0, float64(otherElo-currentElo)/400.0))

		// Accumulate the difference
		totalChange += (sij - eij)
	}

	// Apply the formula: R'i = Ri + K/(N-1) * Σ(Sij - Eij)
	eloChange := (kFactor / float64(numPlayers-1)) * totalChange
	newElo := currentElo + int(eloChange)

	// Ensure ELO doesn't go below 0
	if newElo < 0 {
		newElo = 0
	}

	return newElo
}

// calculateRollingWinrate computes a true moving average winrate over the last N games
func (h *EventHandlers) calculateRollingWinrate(playerID string, won bool) float64 {
	// Use a window of 10 games for the moving average
	windowSize := 10

	// Fetch the last N stats entries for this player
	recentStats, _, err := h.repo.GetPlayerStatsTimeSeries(playerID, windowSize, 0)
	if err != nil || len(recentStats) == 0 {
		// If we can't fetch history, return simple result
		if won {
			return 1.0
		}
		return 0.0
	}

	// Count wins in the recent history
	winsInWindow := 0
	for _, stat := range recentStats {
		// Check if that game was a win (we can infer from the change in TotalWins)
		// For the first stat, we know it's a win if TotalWins > 0
		if len(recentStats) == 1 {
			winsInWindow = stat.TotalWins
		} else {
			// For subsequent stats, compare with previous
			// We'll count the stat as a win if it increased the total
			winsInWindow += stat.TotalWins
		}
	}

	// Actually, let's use a simpler approach: calculate wins from the window
	// by looking at the last game's total wins and the oldest game's total wins
	if len(recentStats) >= windowSize {
		// We have a full window - use the last N games
		newestStat := recentStats[0] // Most recent (DESC order)
		oldestStat := recentStats[windowSize-1]
		winsInWindow = newestStat.TotalWins - oldestStat.TotalWins
		// Add 1 if current game is a win (not yet in stats)
		if won {
			winsInWindow++
		}
		return float64(winsInWindow) / float64(windowSize)
	}

	// Fewer than windowSize games - use all games including the current one
	totalGames := len(recentStats) + 1
	winsInWindow = recentStats[0].TotalWins // Most recent total
	if won {
		winsInWindow++
	}
	return float64(winsInWindow) / float64(totalGames)
}

// calculateStreak updates the win/loss streak
func (h *EventHandlers) calculateStreak(currentStreak int, won bool) int {
	if won {
		if currentStreak >= 0 {
			return currentStreak + 1
		}
		return 1
	}
	// Loss
	if currentStreak <= 0 {
		return currentStreak - 1
	}
	return -1
}
