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
		Image:      game.Image,
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
		sourcePlayer = event.SourceRanking.Player.Name
		sourceCommander = event.SourceRanking.Deck.Commander
	}

	if event.TargetRanking != nil {
		targetPlayer = event.TargetRanking.Player.Name
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
			Deck: Deck{
				Commander:    rank.Deck.Commander,
				Crop:         rank.Deck.Crop,
				SecondaryImg: rank.Deck.SecondaryImage,
				Image:        rank.Deck.Image,
			},
			Player: Player{
				ID:   rank.Player.ID,
				Name: rank.Player.Name,
			},
		}
	}
	return result
}

func convertPlayerToDto(player *repository.Player) Player {
	result := Player{
		ID:   player.ID,
		Name: player.Name,
	}

	// Calculate winrate and game statistics
	totalGames := len(player.Games)
	wins := 0
	deckMap := make(map[string]DeckWithCount)
	coPlayerMap := make(map[uint]PlayerWithCount)
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
			if ranking.PlayerID == player.ID {
				// Count wins (position 1)
				if ranking.Position == 1 {
					wins++
				}

				// Collect unique decks
				deckKey := ranking.Deck.Commander
				if _, exists := deckMap[deckKey]; !exists {
					deckMap[deckKey] = DeckWithCount{Deck: Deck{
						Commander:    ranking.Deck.Commander,
						Crop:         ranking.Deck.Crop,
						SecondaryImg: ranking.Deck.SecondaryImage,
						Image:        ranking.Deck.Image,
					},
					}
				} else {
					deckMap[deckKey] = DeckWithCount{
						Deck:  deckMap[deckKey].Deck,
						Count: deckMap[deckKey].Count + 1,
					}
				}
				break
			}
		}

		// Collect co-players
		for _, ranking := range game.Rankings {
			if ranking.PlayerID != player.ID {
				if coPlayer, exists := coPlayerMap[ranking.PlayerID]; exists {
					coPlayer.Count++
					coPlayerMap[ranking.PlayerID] = coPlayer
				} else {
					coPlayerMap[ranking.PlayerID] = PlayerWithCount{
						Player: Player{
							ID:   ranking.Player.ID,
							Name: ranking.Player.Name,
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
