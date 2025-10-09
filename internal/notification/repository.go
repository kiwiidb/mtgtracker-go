package notification

import (
	"errors"
	"fmt"
	"log"
	"mtgtracker/internal/mtgtracker"
	"mtgtracker/internal/repository"

	"gorm.io/gorm"
)

type NotificationAction string

const (
	ActionDeleteRanking     NotificationAction = "delete_ranking"
	ActionViewGame          NotificationAction = "view_game"
	ActionAddImageGameEvent NotificationAction = "add_image_game_event"
)

type Notification struct {
	gorm.Model
	UserID           string               `gorm:"not null" json:"user_id"`
	Title            string               `gorm:"not null" json:"title"`
	Body             string               `gorm:"not null" json:"body"`
	Type             string               `gorm:"not null" json:"type"`
	Actions          []NotificationAction `gorm:"serializer:json" json:"actions"`
	Read             bool                 `gorm:"default:false" json:"read"`
	GameID           *uint                `json:"game_id,omitempty"`
	ReferredPlayerID *string              `json:"referred_player_id,omitempty"`
	PlayerRankingID  *uint                `json:"player_ranking_id,omitempty"`

	Player         mtgtracker.Player   `gorm:"foreignKey:UserID;references:FirebaseID" json:"user"`
	Game           *repository.Game    `gorm:"foreignKey:GameID;references:ID" json:"game,omitempty"`
	ReferredPlayer *repository.Player  `gorm:"foreignKey:ReferredPlayerID;references:FirebaseID" json:"referred_player,omitempty"`
	PlayerRanking  *repository.Ranking `gorm:"foreignKey:PlayerRankingID;references:ID" json:"player_ranking,omitempty"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) GetNotifications(userID string, readFilter *bool, limit, offset int) ([]Notification, int64, error) {
	var notifications []Notification
	var total int64

	query := r.DB.Model(&Notification{}).Where("user_id = ?", userID)

	// Apply read filter if provided
	if readFilter != nil {
		query = query.Where("read = ?", *readFilter)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.DB.Where("user_id = ?", userID).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if readFilter != nil {
				return db.Where("read = ?", *readFilter)
			}
			return db
		}).
		Preload("Player").
		Preload("Game.Rankings.Player").
		Preload("Game.Rankings.Deck").
		Preload("Game.GameEvents").
		Preload("Game.Creator").
		Preload("ReferredPlayer").
		Preload("PlayerRanking").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&notifications).Error
	if err != nil {
		return nil, 0, err
	}
	for _, n := range notifications {
		log.Println(len(n.Game.Rankings))
		for _, rk := range n.Game.Rankings {
			log.Println(rk.ID)

		}
	}
	return notifications, total, nil
}

func (r *Repository) MarkNotificationAsRead(notificationID uint, userID string) error {
	result := r.DB.Model(&Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("read", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("notification not found or access denied")
	}
	return nil
}

func (r *Repository) createGameCreatedNotifications(game *mtgtracker.Game, creator *mtgtracker.Player) error {
	// Create notifications for all players in the game
	for _, ranking := range game.Rankings {
		if ranking.PlayerID != nil {
			// Get commander name from either referenced deck or embedded deck
			commanderName := ""
			commanderName = ranking.Deck.Commander

			notification := Notification{
				UserID:           *ranking.PlayerID,
				ReferredPlayerID: game.CreatorID,
				Title:            fmt.Sprintf("%s started a game", creator.Name),
				Body:             fmt.Sprintf("You're playing %s", commanderName),
				Type:             "game_created",
				Actions:          []NotificationAction{ActionViewGame, ActionDeleteRanking, ActionAddImageGameEvent},
				Read:             false,
				GameID:           &game.ID,
				PlayerRankingID:  &ranking.ID,
			}

			if err := r.DB.Create(&notification).Error; err != nil {
				log.Printf("Failed to create notification for player %s: %v", *ranking.PlayerID, err)
				// Continue creating notifications for other players even if one fails
			}
		}
	}
	return nil
}

func (r *Repository) CreateFinishedGameNotifications(game *mtgtracker.Game) error {
	// Delete all game_created notifications for this game
	if err := r.DB.Where("game_id = ? AND type = ?", game.ID, "game_created").Delete(&Notification{}).Error; err != nil {
		log.Printf("Failed to delete game_created notifications for game %d: %v", game.ID, err)
		// Don't fail the entire operation if deletion fails
	}

	// Get all player names for the body
	playerNames := make([]string, 0, len(game.Rankings))
	for _, ranking := range game.Rankings {
		if ranking.Player != nil {
			playerNames = append(playerNames, ranking.Player.Name)
		} else {
			playerNames = append(playerNames, "guest")
		}
	}

	// Create notifications for all players
	for _, ranking := range game.Rankings {
		if ranking.PlayerID != nil {
			// Build list of other players (excluding current player)
			otherPlayers := make([]string, 0, len(playerNames)-1)
			for _, name := range playerNames {
				if ranking.Player != nil && name != ranking.Player.Name {
					otherPlayers = append(otherPlayers, name)
				}
			}

			var title, notificationType string
			if ranking.Position == 1 {
				title = "🎉 You won the game!"
				notificationType = "game_finished_won"
			} else {
				title = "Game finished"
				notificationType = "game_finished"
			}

			body := "Opponents "
			if len(otherPlayers) > 0 {
				body += otherPlayers[0]
				for i := 1; i < len(otherPlayers); i++ {
					if i == len(otherPlayers)-1 {
						body += " and " + otherPlayers[i]
					} else {
						body += ", " + otherPlayers[i]
					}
				}
			}

			notification := Notification{
				UserID:           *ranking.PlayerID,
				ReferredPlayerID: game.CreatorID,
				Title:            title,
				Body:             body,
				Type:             notificationType,
				Actions:          []NotificationAction{ActionViewGame, ActionDeleteRanking},
				Read:             false,
				GameID:           &game.ID,
				PlayerRankingID:  &ranking.ID,
			}

			if err := r.DB.Create(&notification).Error; err != nil {
				log.Printf("Failed to create finished game notification for player %s: %v", *ranking.PlayerID, err)
				// Continue creating notifications for other players even if one fails
			}
		}
	}
	return nil
}
