package mtgtracker

import "mtgtracker/internal/repository"

func convertGameToDto(game *repository.Game) GameDto {
	result := GameDto{
		ID:         game.ID,
		Duration:   game.Duration,
		Date:       game.Date,
		Comments:   game.Comments,
		Image:      game.Image,
		Finished:   game.Finished,
		Rankings:   convertRankingsToDto(game.Rankings),
		GameEvents: make([]GameEventDto, len(game.GameEvents)),
	}
	for i, event := range game.GameEvents {
		result.GameEvents[i] = convertGameEvent(&event, "")
	}
	return result
}

func convertGameEvent(event *repository.GameEvent, uploadUrl string) GameEventDto {
	var sourcePlayer, targetPlayer string
	var sourceCommander, targetCommander string
	var sourceCropImg, targetCropImg string

	if event.SourceRanking != nil {
		sourcePlayer = event.SourceRanking.Player.Name
		sourceCommander = event.SourceRanking.Deck.Commander
		sourceCropImg = event.SourceRanking.Deck.Crop
	}

	if event.TargetRanking != nil {
		targetPlayer = event.TargetRanking.Player.Name
		targetCommander = event.TargetRanking.Deck.Commander
		targetCropImg = event.TargetRanking.Deck.Crop
	}

	return GameEventDto{
		GameID:                 event.GameID,
		EventType:              event.EventType,
		DamageDelta:            event.DamageDelta,
		CreatedAt:              event.CreatedAt,
		TargetLifeTotalAfter:   event.TargetLifeTotalAfter,
		SourcePlayer:           sourcePlayer,
		TargetPlayer:           targetPlayer,
		SourceCommander:        sourceCommander,
		TargetCommander:        targetCommander,
		SourceCommanderCropImg: sourceCropImg,
		TargetCommanderCropImg: targetCropImg,
		ImageUrl:               event.ImageUrl,
		UploadImageUrl:         uploadUrl,
	}
}

func convertRankingsToDto(rankings []repository.Ranking) []Ranking {
	result := make([]Ranking, len(rankings))
	for i, rank := range rankings {
		result[i] = Ranking{
			ID:             rank.ID,
			PlayerID:       rank.PlayerID,
			Position:       rank.Position,
			CouldHaveWon:   rank.CouldHaveWon,
			EarlySolRing:   rank.EarlySolRing,
			StartingPlayer: rank.StartingPlayer,
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

func convertPlayerToDto(player *repository.Player) PlayerDto {
	result := PlayerDto{
		ID:   player.ID,
		Name: player.Name,
	}

	// Calculate winrate and game statistics
	totalGames := len(player.Games)
	wins := 0
	deckMap := make(map[string]Deck)
	coPlayerMap := make(map[uint]Player)
	games := make([]GameDto, len(player.Games))

	for i, game := range player.Games {
		// Convert game to DTO
		games[i] = convertGameToDto(&game)

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
					deckMap[deckKey] = Deck{
						Commander:    ranking.Deck.Commander,
						Crop:         ranking.Deck.Crop,
						SecondaryImg: ranking.Deck.SecondaryImage,
						Image:        ranking.Deck.Image,
					}
				}
				break
			}
		}

		// Collect co-players
		for _, ranking := range game.Rankings {
			if ranking.PlayerID != player.ID {
				coPlayerMap[ranking.Player.ID] = Player{
					ID:   ranking.Player.ID,
					Name: ranking.Player.Name,
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
	decks := make([]Deck, 0, len(deckMap))
	for _, deck := range deckMap {
		decks = append(decks, deck)
	}

	coPlayers := make([]Player, 0, len(coPlayerMap))
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
