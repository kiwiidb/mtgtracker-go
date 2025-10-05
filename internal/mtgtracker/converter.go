package mtgtracker

import (
	"mtgtracker/internal/repository"
	"path/filepath"
	"sort"
)

func convertGameToDto(game *repository.Game) Game {
	result := Game{
		ID:         game.ID,
		CreatorID:  game.CreatorID,
		Duration:   game.Duration,
		Date:       game.Date,
		Comments:   game.Comments,
		Finished:   game.Finished,
		Rankings:   convertRankingsToDto(game.Rankings),
		GameEvents: make([]GameEvent, len(game.GameEvents)),
	}

	// Include creator information if available
	if game.Creator != nil && game.Creator.FirebaseID != "" {
		creator := convertPlayerToDto(game.Creator)
		result.Creator = &creator
	}

	for i, event := range game.GameEvents {
		result.GameEvents[i] = convertGameEvent(&event, "")
	}
	return result
}

func convertDeck(deck Deck) repository.Deck {
	return repository.Deck{
		Commander:      deck.Commander,
		Image:          deck.Image,
		SecondaryImage: deck.SecondaryImg,
		Crop:           deck.Crop,
	}
}

func convertGameEvent(event *repository.GameEvent, uploadUrl string) GameEvent {
	var sourcePlayer, targetPlayer string
	var sourceCommander, targetCommander string

	if event.SourceRanking != nil {
		if event.SourceRanking.Player != nil {
			sourcePlayer = event.SourceRanking.Player.Name
		}
		// Get commander from either referenced deck or embedded deck
		if event.SourceRanking.Deck != nil {
			sourceCommander = event.SourceRanking.Deck.Commander
		} else {
			sourceCommander = event.SourceRanking.DeckEmbedded.Commander
		}
	}

	if event.TargetRanking != nil {
		if event.TargetRanking.Player != nil {
			targetPlayer = event.TargetRanking.Player.Name
		}
		// Get commander from either referenced deck or embedded deck
		if event.TargetRanking.Deck != nil {
			targetCommander = event.TargetRanking.Deck.Commander
		} else {
			targetCommander = event.TargetRanking.DeckEmbedded.Commander
		}
	}

	return GameEvent{
		GameID:               event.GameID,
		EventType:            event.EventType,
		DamageDelta:          event.DamageDelta,
		CreatedAt:            event.CreatedAt,
		TargetLifeTotalAfter: event.TargetLifeTotalAfter,
		SourcePlayer:         sourcePlayer,
		TargetPlayer:         targetPlayer,
		SourceCommander:      sourceCommander,
		TargetCommander:      targetCommander,
		ImageUrl:             event.ImageUrl,
		UploadImageUrl:       uploadUrl,
	}
}

func convertRankingsToDto(rankings []repository.Ranking) []Ranking {
	result := make([]Ranking, len(rankings))
	for i, rank := range rankings {
		// Determine which deck to use: referenced deck or embedded deck
		var deckData Deck
		if rank.Deck != nil {
			// Use referenced deck from player's deck collection
			deckData = Deck{
				Commander:    rank.Deck.Commander,
				Crop:         rank.Deck.Crop,
				SecondaryImg: rank.Deck.SecondaryImage,
				Image:        rank.Deck.Image,
			}
		} else {
			// Use embedded deck data
			deckData = Deck{
				Commander:    rank.DeckEmbedded.Commander,
				Crop:         rank.DeckEmbedded.Crop,
				SecondaryImg: rank.DeckEmbedded.SecondaryImage,
				Image:        rank.DeckEmbedded.Image,
			}
		}

		result[i] = Ranking{
			ID:       rank.ID,
			PlayerID: rank.PlayerID,
			Position: rank.Position,
			Deck:     deckData,
			Player: func() *Player {
				if rank.Player != nil {
					return &Player{
						ID:              rank.Player.FirebaseID,
						Name:            rank.Player.Name,
						ProfileImageURL: rank.Player.Image,
					}
				}
				return nil
			}(),
		}
	}
	return result
}

func convertPlayerToDto(player *repository.Player) Player {
	result := Player{
		ID:               player.FirebaseID,
		Name:             player.Name,
		ProfileImageURL:  player.Image,
		MoxfieldUsername: player.MoxfieldUsername,
	}

	// Calculate winrate and game statistics
	totalGames := len(player.Games)
	wins := 0
	deckMap := make(map[string]DeckWithCount)
	coPlayerMap := make(map[string]PlayerWithCount)
	games := make([]Game, len(player.Games))

	for i, game := range player.Games {
		// Convert game to DTO
		games[i] = convertGameToDto(&game)

		// Normally you will only have 1 game in progress at a time
		if !game.Finished {
			result.CurrentGame = &games[i]
			continue
		}

		// Find this player's ranking in the game
		for _, ranking := range game.Rankings {
			if ranking.PlayerID != nil && *ranking.PlayerID == player.FirebaseID {
				deckKey := ranking.Deck.Commander
				if _, exists := deckMap[deckKey]; !exists {
					deckMap[deckKey] = DeckWithCount{Deck: Deck{
						Commander:    ranking.Deck.Commander,
						Crop:         ranking.Deck.Crop,
						SecondaryImg: ranking.Deck.SecondaryImage,
						Image:        ranking.Deck.Image,
					},
						Count: 1,
					}

				} else {
					deckMap[deckKey] = DeckWithCount{
						Deck:  deckMap[deckKey].Deck,
						Count: deckMap[deckKey].Count + 1,
						Wins:  deckMap[deckKey].Wins, // Keep the wins count intact
					}
				}
				if ranking.Position == 1 {
					wins++
					entry := deckMap[deckKey]
					entry.Wins++
					deckMap[deckKey] = entry
				}
			}
		}

		// Collect co-players
		for _, ranking := range game.Rankings {
			if ranking.PlayerID != nil && *ranking.PlayerID != player.FirebaseID {
				if coPlayer, exists := coPlayerMap[*ranking.PlayerID]; exists {
					coPlayer.Count++
					coPlayerMap[*ranking.PlayerID] = coPlayer
				} else {
					coPlayerMap[*ranking.PlayerID] = PlayerWithCount{
						Player: Player{
							ID: func() string {
								if ranking.Player != nil {
									return ranking.Player.FirebaseID
								}
								if ranking.PlayerID != nil {
									return *ranking.PlayerID
								}
								return ""
							}(),
							Name: func() string {
								if ranking.Player != nil {
									return ranking.Player.Name
								}
								return ""
							}(),
							ProfileImageURL: func() string {
								if ranking.Player != nil {
									return ranking.Player.Image
								}
								return ""
							}(),
						},
						Count: 1,
					}
				}
			}
		}
	}

	// Calculate winrate
	var winrate float64
	if totalGames > 0 {
		winrate = float64(wins) / float64(totalGames) * 100
	}

	// Convert maps to slices
	decks := make([]DeckWithCount, 0, len(deckMap))
	for _, deck := range deckMap {
		decks = append(decks, deck)
	}
	// sort decks by count descending
	sort.Slice(decks, func(i, j int) bool {
		return decks[i].Count > decks[j].Count
	})

	coPlayers := make([]PlayerWithCount, 0, len(coPlayerMap))
	for _, coPlayer := range coPlayerMap {
		coPlayers = append(coPlayers, coPlayer)
	}
	//sort coPlayers by count descending
	sort.Slice(coPlayers, func(i, j int) bool {
		return coPlayers[i].Count > coPlayers[j].Count
	})

	result.WinrateAllTime = winrate
	result.NumberofGamesAllTime = totalGames
	result.DecksAllTime = decks
	result.CoPlayersAllTime = coPlayers
	result.Games = games

	return result
}

// convertPlayerToDtoSimple converts a player for notification context without game statistics
func convertPlayerToDtoSimple(player *repository.Player) Player {
	return Player{
		ID:                   player.FirebaseID,
		Name:                 player.Name,
		ProfileImageURL:      player.Image,
		MoxfieldUsername:     player.MoxfieldUsername,
		WinrateAllTime:       0,
		NumberofGamesAllTime: 0,
		DecksAllTime:         []DeckWithCount{},
		CoPlayersAllTime:     []PlayerWithCount{},
		Games:                []Game{},
		CurrentGame:          nil,
	}
}

func convertNotificationToDto(notification *repository.Notification) Notification {
	result := Notification{
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
		game := convertGameToDto(notification.Game)
		result.Game = &game
	}

	if notification.ReferredPlayer != nil {
		player := convertPlayerToDtoSimple(notification.ReferredPlayer)
		result.ReferredPlayer = &player
	}

	return result
}

func convertActionsToDto(actions []repository.NotificationAction) []NotificationAction {
	result := make([]NotificationAction, len(actions))
	for i, action := range actions {
		result[i] = NotificationAction(action)
	}
	return result
}

func convertDeckToDto(deck *repository.Deck) Deck {
	return Deck{
		Commander:    deck.Commander,
		Crop:         deck.Crop,
		SecondaryImg: deck.SecondaryImage,
		Image:        deck.Image,
	}
}

func getImgContentType(s string) string {
	switch filepath.Ext(s) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}
