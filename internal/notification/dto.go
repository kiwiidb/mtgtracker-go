package notification

import (
	"mtgtracker/internal/mtgtracker"
	"time"
)

type NotificationResponse struct {
	ID               uint                 `json:"id"`
	Title            string               `json:"title"`
	Body             string               `json:"body"`
	Type             string               `json:"type"`
	Actions          []NotificationAction `json:"actions"`
	Read             bool                 `json:"read"`
	CreatedAt        time.Time            `json:"created_at"`
	GameID           *uint                `json:"game_id,omitempty"`
	ReferredPlayerID *string              `json:"referred_player_id,omitempty"`
	PlayerRankingID  *uint                `json:"player_ranking_id,omitempty"`
	Game             *mtgtracker.Game     `json:"game,omitempty"`
	ReferredPlayer   *mtgtracker.Player   `json:"referred_player,omitempty"`
}

type NotificationResponseAction string

const (
	ActionResponseDeleteRanking     NotificationResponseAction = "delete_ranking"
	ActionResponseViewGame          NotificationResponseAction = "view_game"
	ActionResponseAddImageGameEvent NotificationResponseAction = "add_image_game_event"
)

func (s *Service) convertNotificationToDto(notification *Notification) NotificationResponse {
	result := NotificationResponse{
		ID:               notification.ID,
		Title:            notification.Title,
		Body:             notification.Body,
		Type:             notification.Type,
		Actions:          convertActionsToDto(notification.Actions),
		Read:             notification.Read,
		CreatedAt:        notification.CreatedAt,
		GameID:           notification.GameID,
		ReferredPlayerID: notification.ReferredPlayerID,
		PlayerRankingID:  notification.PlayerRankingID,
	}

	if notification.Game != nil {
		game := s.coreService.ConvertGameToDto(notification.Game, false)
		result.Game = &game
	}

	if notification.ReferredPlayer != nil {
		player := s.coreService.ConvertPlayerToResponse(notification.ReferredPlayer)
		result.ReferredPlayer = &player
	}

	return result
}

func convertActionsToDto(actions []NotificationAction) []NotificationAction {
	result := make([]NotificationAction, len(actions))
	for i, action := range actions {
		result[i] = NotificationAction(action)
	}
	return result
}
