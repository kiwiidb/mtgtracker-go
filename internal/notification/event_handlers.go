package notification

import (
	"fmt"
	"log"
	"mtgtracker/internal/events"
)

type PushService interface {
	SendNotification(playerID, title, body string, data map[string]string) error
}

// EventHandlers manages event subscriptions for the notification package
type EventHandlers struct {
	repo        *Repository
	coreService CoreService
	pushService PushService
}

// NewEventHandlers creates a new event handler instance
func NewEventHandlers(repo *Repository, coreService CoreService, pushService PushService) *EventHandlers {
	return &EventHandlers{
		repo:        repo,
		coreService: coreService,
		pushService: pushService,
	}
}

// RegisterHandlers subscribes to all relevant events
func (h *EventHandlers) RegisterHandlers(bus *events.EventBus) {
	bus.Subscribe("game.created", h.HandleGameCreated)
	bus.Subscribe("game.finished", h.HandleGameFinished)
	bus.Subscribe("game.deleted", h.HandleGameDeleted)
	log.Println("Notification event handlers registered")
}

// HandleGameCreated processes game created events
func (h *EventHandlers) HandleGameCreated(event events.Event) error {
	e, ok := event.(events.GameCreatedEvent)
	if !ok {
		log.Printf("Invalid event type for game.created: %T", event)
		return nil
	}

	log.Printf("Processing game.created event for game %d", e.GameID)

	// Fetch game and creator data via service interface
	game, err := h.coreService.GetGameByID(e.GameID)
	if err != nil {
		log.Printf("Failed to fetch game %d: %v", e.GameID, err)
		return err
	}

	creator, err := h.coreService.GetPlayerByFirebaseID(e.CreatorID)
	if err != nil {
		log.Printf("Failed to fetch creator %s: %v", e.CreatorID, err)
		return err
	}

	// Create in-app notifications
	err = h.repo.CreateGameNotifications(game, creator)
	if err != nil {
		return err
	}

	// Send push notifications to all players except creator
	for _, ranking := range game.Rankings {
		if ranking.PlayerID != nil && *ranking.PlayerID != e.CreatorID {
			err := h.pushService.SendNotification(
				*ranking.PlayerID,
				"New Game",
				fmt.Sprintf("%s created a game with you", creator.Name),
				map[string]string{
					"type":    "game_created",
					"game_id": fmt.Sprint(e.GameID),
				},
			)
			if err != nil {
				log.Printf("Failed to send push notification to %s: %v", *ranking.PlayerID, err)
				// Continue processing other players
			}
		}
	}

	return nil
}

// HandleGameFinished processes game finished events
func (h *EventHandlers) HandleGameFinished(event events.Event) error {
	e, ok := event.(events.GameFinishedEvent)
	if !ok {
		log.Printf("Invalid event type for game.finished: %T", event)
		return nil
	}

	log.Printf("Processing game.finished event for game %d", e.GameID)

	// Fetch game data
	game, err := h.coreService.GetGameByID(e.GameID)
	if err != nil {
		log.Printf("Failed to fetch game %d: %v", e.GameID, err)
		return err
	}

	// Create in-app notifications
	err = h.repo.CreateGameFinishedNotifications(game)
	if err != nil {
		return err
	}

	// Send push notifications to all players
	// Find the winner
	var winnerName string
	for _, ranking := range game.Rankings {
		if ranking.Position == 1 && ranking.Player != nil {
			winnerName = ranking.Player.Name
			break
		}
	}

	for _, ranking := range game.Rankings {
		if ranking.PlayerID != nil {
			var body string
			if ranking.Position == 1 {
				body = "Congratulations! You won the game!"
			} else if winnerName != "" {
				body = fmt.Sprintf("Game finished. %s won!", winnerName)
			} else {
				body = "Game finished!"
			}

			err := h.pushService.SendNotification(
				*ranking.PlayerID,
				"Game Finished",
				body,
				map[string]string{
					"type":     "game_finished",
					"game_id":  fmt.Sprint(e.GameID),
					"position": fmt.Sprint(ranking.Position),
				},
			)
			if err != nil {
				log.Printf("Failed to send push notification to %s: %v", *ranking.PlayerID, err)
				// Continue processing other players
			}
		}
	}

	return nil
}

// HandleGameDeleted processes game deleted events
func (h *EventHandlers) HandleGameDeleted(event events.Event) error {
	e, ok := event.(events.GameDeletedEvent)
	if !ok {
		log.Printf("Invalid event type for game.deleted: %T", event)
		return nil
	}

	log.Printf("Processing game.deleted event for notifications (game %d)", e.GameID)

	// Delete all notifications related to this game
	return h.repo.DeleteNotificationsByGameID(e.GameID)
}
