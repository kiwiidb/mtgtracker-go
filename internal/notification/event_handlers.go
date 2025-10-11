package notification

import (
	"log"
	"mtgtracker/internal/events"
)

// EventHandlers manages event subscriptions for the notification package
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

	return h.repo.CreateGameNotifications(game, creator)
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

	return h.repo.CreateGameFinishedNotifications(game)
}

// HandleGameDeleted processes game deleted events
func (h *EventHandlers) HandleGameDeleted(event events.Event) error {
	e, ok := event.(events.GameDeletedEvent)
	if !ok {
		log.Printf("Invalid event type for game.deleted: %T", event)
		return nil
	}

	log.Printf("Processing game.deleted event for game %d", e.GameID)

	// Delete all notifications related to this game
	return h.repo.DeleteNotificationsByGameID(e.GameID)
}
