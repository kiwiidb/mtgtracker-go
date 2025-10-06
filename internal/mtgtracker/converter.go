package mtgtracker

import (
	"mtgtracker/internal/repository"
	"path/filepath"
	"sort"
	"time"
)

func convertGameToDto(game *repository.Game) Game {
	result := Game{
		ID:         game.ID,
		CreatorID:  game.CreatorID,
		Duration:   game.Duration,
		Date:       game.Date,
		Comments:   game.Comments,
		Finished:   game.Finished,
		Rankings:   convertRankingsWithLifeTotal(game.Rankings, game.GameEvents),
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

func convertSimpleDeck(deck Deck) repository.SimpleDeck {
	return repository.SimpleDeck{
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

func convertRankingsWithLifeTotal(rankings []repository.Ranking, gameEvents []repository.GameEvent) []Ranking {
	result := make([]Ranking, len(rankings))

	// Build a map of ranking ID to most recent life total event
	lifeTotalMap := make(map[uint]*repository.GameEvent)
	for i := range gameEvents {
		event := &gameEvents[i]
		if event.TargetRankingID != nil {
			// Keep the most recent event (events are assumed to be sorted by CreatedAt)
			// If not sorted, we'll take the last one which should be most recent
			lifeTotalMap[*event.TargetRankingID] = event
		}
	}

	for i, rank := range rankings {
		// Determine which deck to use: referenced deck or embedded deck
		// Get last life total from most recent event for this ranking
		var lastLifeTotal *int
		var lastLifeTotalTimestamp *time.Time
		if event, exists := lifeTotalMap[rank.ID]; exists {
			lastLifeTotal = &event.TargetLifeTotalAfter
			timestamp := event.CreatedAt
			lastLifeTotalTimestamp = &timestamp
		}
		var deckData Deck
		if rank.Deck != nil {
			// Use referenced deck from player's deck collection
			deckData = Deck{
				Commander:    rank.Deck.Commander,
				Colors:       rank.Deck.Colors,
				Crop:         rank.Deck.Crop,
				SecondaryImg: rank.Deck.SecondaryImage,
				Image:        rank.Deck.Image,
			}
		} else {
			// Use embedded deck data
			deckData = Deck{
				Commander:    rank.DeckEmbedded.Commander,
				Colors:       nil, // Embedded decks don't store colors
				Crop:         rank.DeckEmbedded.Crop,
				SecondaryImg: rank.DeckEmbedded.SecondaryImage,
				Image:        rank.DeckEmbedded.Image,
			}
		}

		result[i] = Ranking{
			ID:                     rank.ID,
			PlayerID:               rank.PlayerID,
			Position:               rank.Position,
			Deck:                   deckData,
			LastLifeTotal:          lastLifeTotal,
			LastLifeTotalTimestamp: lastLifeTotalTimestamp,
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

		// Find this player's ranking in the game to count wins
		for _, ranking := range game.Rankings {
			if ranking.PlayerID != nil && *ranking.PlayerID == player.FirebaseID {
				if ranking.Position == 1 {
					wins++
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

	// Convert player.Decks to DeckWithCount
	decks := make([]DeckWithCount, 0, len(player.Decks))
	for _, deck := range player.Decks {
		decks = append(decks, DeckWithCount{
			Deck: Deck{
				ID:           &deck.ID,
				Commander:    deck.Commander,
				Colors:       deck.Colors,
				Crop:         deck.Crop,
				SecondaryImg: deck.SecondaryImage,
				Image:        deck.Image,
			},
			Count: deck.GameCount,
			Wins:  deck.WinCount,
		})
	}
	// sort decks by count descending
	sort.Slice(decks, func(i, j int) bool {
		return decks[i].Count > decks[j].Count
	})

	// Calculate top 2 most played colors from decks
	colors := calculateTopColors(decks, 2)

	coPlayers := make([]PlayerWithCount, 0, len(coPlayerMap))
	for _, coPlayer := range coPlayerMap {
		coPlayers = append(coPlayers, coPlayer)
	}
	//sort coPlayers by count descending
	sort.Slice(coPlayers, func(i, j int) bool {
		return coPlayers[i].Count > coPlayers[j].Count
	})

	result.Colors = colors
	result.WinrateAllTime = winrate
	result.NumberofGamesAllTime = totalGames
	result.DecksAllTime = decks
	result.CoPlayersAllTime = coPlayers
	result.Games = games

	return result
}

// calculateTopColors calculates the most played colors from decks
func calculateTopColors(decks []DeckWithCount, topN int) []string {
	// Count color occurrences weighted by game count
	colorCounts := make(map[string]int)

	for _, deckWithCount := range decks {
		for _, color := range deckWithCount.Deck.Colors {
			colorCounts[color] += deckWithCount.Count
		}
	}

	// Convert to slice for sorting
	type colorCount struct {
		color string
		count int
	}
	colorSlice := make([]colorCount, 0, len(colorCounts))
	for color, count := range colorCounts {
		colorSlice = append(colorSlice, colorCount{color: color, count: count})
	}

	// Sort by count descending
	sort.Slice(colorSlice, func(i, j int) bool {
		return colorSlice[i].count > colorSlice[j].count
	})

	// Get top N colors
	result := make([]string, 0, topN)
	for i := 0; i < topN && i < len(colorSlice); i++ {
		result = append(result, colorSlice[i].color)
	}

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
		ID:           &deck.ID,
		Commander:    deck.Commander,
		Colors:       deck.Colors,
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
