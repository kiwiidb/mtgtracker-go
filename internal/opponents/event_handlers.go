package opponents

import (
	"log"
	"mtgtracker/internal/core"
	"mtgtracker/internal/events"
)

// EventHandlers manages event subscriptions for the opponents package
type EventHandlers struct {
	repo        *Repository
	coreService CoreService
}

type CoreService interface {
	GetGameByID(gameID uint) (*core.Game, error)
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
	bus.Subscribe("game.created", h.HandleGameCreated)
	bus.Subscribe("game.deleted", h.HandleGameDeleted)
	bus.Subscribe("ranking.deleted", h.HandleRankingDeleted)
	log.Println("Opponent event handlers registered")
}

// HandleGameCreated processes game created events
// Creates or updates follow relationships for all player pairs in the game
func (h *EventHandlers) HandleGameCreated(event events.Event) error {
	e, ok := event.(events.GameCreatedEvent)
	if !ok {
		log.Printf("Invalid event type for game.created: %T", event)
		return nil
	}

	log.Printf("Processing game.created event for opponents (game %d)", e.GameID)

	// Fetch game data to get rankings
	game, err := h.coreService.GetGameByID(e.GameID)
	if err != nil {
		log.Printf("Failed to fetch game %d: %v", e.GameID, err)
		return err
	}

	// Extract player IDs from rankings (excluding guest players)
	playerIDs := make([]string, 0, len(game.Rankings))
	for _, ranking := range game.Rankings {
		if ranking.PlayerID != nil {
			playerIDs = append(playerIDs, *ranking.PlayerID)
		}
	}

	// Create/update opponents for all unique player pairs
	return h.updateOpponentsForPlayerPairs(playerIDs, true)
}

// HandleGameDeleted processes game deleted events
// Decrements follow counts for all player pairs in the game
func (h *EventHandlers) HandleGameDeleted(event events.Event) error {
	e, ok := event.(events.GameDeletedEvent)
	if !ok {
		log.Printf("Invalid event type for game.deleted: %T", event)
		return nil
	}

	log.Printf("Processing game.deleted event for opponents (game %d)", e.GameID)

	// Use player IDs from the event (game is already deleted)
	// Decrement opponents for all unique player pairs
	return h.updateOpponentsForPlayerPairs(e.PlayerIDs, false)
}

// HandleRankingDeleted processes ranking deleted events
// Decrements follow counts between the deleted player and all other players in the game
func (h *EventHandlers) HandleRankingDeleted(event events.Event) error {
	e, ok := event.(events.RankingDeletedEvent)
	if !ok {
		log.Printf("Invalid event type for ranking.deleted: %T", event)
		return nil
	}

	log.Printf("Processing ranking.deleted event for opponents (ranking %d, game %d)", e.RankingID, e.GameID)

	// Decrement opponents between the removed player and all other players
	for _, otherPlayerID := range e.OtherPlayerIDs {
		err := h.repo.DecrementGameCount(e.PlayerID, otherPlayerID)
		if err != nil {
			log.Printf("Failed to decrement follow count for %s <-> %s: %v", e.PlayerID, otherPlayerID, err)
			// Continue processing other pairs
		}
	}

	return nil
}

// updateOpponentsForPlayerPairs creates or updates follow relationships for all unique pairs
// If increment is true, increments counts; otherwise decrements
func (h *EventHandlers) updateOpponentsForPlayerPairs(playerIDs []string, increment bool) error {
	// Process all unique pairs
	for i := 0; i < len(playerIDs); i++ {
		for j := i + 1; j < len(playerIDs); j++ {
			player1ID := playerIDs[i]
			player2ID := playerIDs[j]

			var err error
			if increment {
				err = h.repo.IncrementGameCount(player1ID, player2ID)
				if err != nil {
					log.Printf("Failed to increment follow count for %s <-> %s: %v", player1ID, player2ID, err)
					// Continue processing other pairs
				}
			} else {
				err = h.repo.DecrementGameCount(player1ID, player2ID)
				if err != nil {
					log.Printf("Failed to decrement follow count for %s <-> %s: %v", player1ID, player2ID, err)
					// Continue processing other pairs
				}
			}
		}
	}

	return nil
}
