package mtgtracker

import (
	"mtgtracker/internal/repository"
	"path/filepath"
)

func convertGameToDto(game *repository.Game) Game {
	result := Game{
		ID:         game.ID,
		Duration:   game.Duration,
		Date:       game.Date,
		Comments:   game.Comments,
		Finished:   game.Finished,
		Rankings:   convertRankingsToDto(game.Rankings),
		GameEvents: make([]GameEvent, len(game.GameEvents)),
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
		sourceCommander = event.SourceRanking.Deck.Commander
	}

	if event.TargetRanking != nil {
		if event.TargetRanking.Player != nil {
			targetPlayer = event.TargetRanking.Player.Name
		}
		targetCommander = event.TargetRanking.Deck.Commander
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
		result[i] = Ranking{
			ID:       rank.ID,
			PlayerID: rank.PlayerID,
			Position: rank.Position,
			Status:   RankingStatus(rank.Status),
			Deck: Deck{
				Commander:    rank.Deck.Commander,
				Crop:         rank.Deck.Crop,
				SecondaryImg: rank.Deck.SecondaryImage,
				Image:        rank.Deck.Image,
			},
			Player: func() *Player {
				if rank.Player != nil {
					return &Player{
						ID:   rank.Player.FirebaseID,
						Name: rank.Player.Name,
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
		ID:   player.FirebaseID,
		Name: player.Name,
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

	coPlayers := make([]PlayerWithCount, 0, len(coPlayerMap))
	for _, coPlayer := range coPlayerMap {
		coPlayers = append(coPlayers, coPlayer)
	}

	result.WinrateAllTime = winrate
	result.NumberofGamesAllTime = totalGames
	result.DecksAllTime = decks
	result.CoPlayersAllTime = coPlayers
	result.Games = games

	return result
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
