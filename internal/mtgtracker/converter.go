package mtgtracker

import (
	"mtgtracker/internal/repository"
	"path/filepath"
	"sort"
	"time"
)

func (svc *Service) ConvertGameToDto(game *repository.Game, addEvents bool) Game {
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
		creator := svc.ConvertPlayerToResponse(game.Creator)
		result.Creator = &creator
	}

	if !addEvents {
		// If not adding events, return early
		return result
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
	var sourceRanking, targetRanking *Ranking

	if event.SourceRanking != nil {
		ranking := convertRankingToDto(event.SourceRanking)
		sourceRanking = &ranking
	}

	if event.TargetRanking != nil {
		ranking := convertRankingToDto(event.TargetRanking)
		targetRanking = &ranking
	}

	return GameEvent{
		GameID:               event.GameID,
		EventType:            event.EventType,
		DamageDelta:          event.DamageDelta,
		CreatedAt:            event.CreatedAt,
		Comment:              event.Comment,
		TargetLifeTotalAfter: event.TargetLifeTotalAfter,
		SourceRanking:        sourceRanking,
		TargetRanking:        targetRanking,
		ImageUrl:             event.ImageUrl,
		UploadImageUrl:       uploadUrl,
	}
}

func convertDeckFromRanking(rank *repository.Ranking) Deck {
	if rank.Deck != nil {
		return Deck{
			Commander:    rank.Deck.Commander,
			Colors:       rank.Deck.Colors,
			Crop:         rank.Deck.Crop,
			SecondaryImg: rank.Deck.SecondaryImage,
			Image:        rank.Deck.Image,
			MoxfieldURL:  rank.Deck.MoxfieldURL,
			Bracket:      rank.Deck.Bracket,
		}
	}
	return Deck{
		Commander:    rank.DeckEmbedded.Commander,
		Colors:       nil,
		Crop:         rank.DeckEmbedded.Crop,
		SecondaryImg: rank.DeckEmbedded.SecondaryImage,
		Image:        rank.DeckEmbedded.Image,
	}
}

func convertRankingToDto(rank *repository.Ranking) Ranking {
	return Ranking{
		ID:       rank.ID,
		PlayerID: rank.PlayerID,
		Position: rank.Position,
		Deck:     convertDeckFromRanking(rank),
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
		// Get last life total from most recent event for this ranking
		var lastLifeTotal *int
		var lastLifeTotalTimestamp *time.Time
		if event, exists := lifeTotalMap[rank.ID]; exists {
			lastLifeTotal = &event.TargetLifeTotalAfter
			timestamp := event.CreatedAt
			lastLifeTotalTimestamp = &timestamp
		}

		result[i] = Ranking{
			ID:                     rank.ID,
			PlayerID:               rank.PlayerID,
			Position:               rank.Position,
			Deck:                   convertDeckFromRanking(&rank),
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

func (svc *Service) ConvertPlayerToResponse(player *repository.Player) Player {
	result := Player{
		ID:               player.FirebaseID,
		Name:             player.Name,
		ProfileImageURL:  player.Image,
		MoxfieldUsername: player.MoxfieldUsername,
	}

	// Calculate winrate and game statistics
	totalGames := len(player.Games)
	wins := 0
	opponentMap := make(map[string]PlayerOpponentWithCount)
	games := make([]Game, len(player.Games))

	for i, game := range player.Games {
		// Convert game to DTO
		games[i] = svc.ConvertGameToDto(&game, false)

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

		// Collect opponents
		for _, ranking := range game.Rankings {
			if ranking.PlayerID != nil && *ranking.PlayerID != player.FirebaseID {
				if opponent, exists := opponentMap[*ranking.PlayerID]; exists {
					opponent.Count++
					opponentMap[*ranking.PlayerID] = opponent
				} else {
					opponentMap[*ranking.PlayerID] = PlayerOpponentWithCount{
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
				MoxfieldURL:  deck.MoxfieldURL,
				Bracket:      deck.Bracket,
				Themes:       deck.Themes,
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

	opponents := make([]PlayerOpponentWithCount, 0, len(opponentMap))
	for _, opponent := range opponentMap {
		opponents = append(opponents, opponent)
	}
	//sort opponents by count descending
	sort.Slice(opponents, func(i, j int) bool {
		return opponents[i].Count > opponents[j].Count
	})

	result.Colors = colors
	result.WinrateAllTime = winrate
	result.NumberofGamesAllTime = totalGames
	result.DecksAllTime = decks
	result.OpponentsAllTime = opponents

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
func ConvertPlayerToDtoSimple(player *repository.Player) Player {
	return Player{
		ID:                   player.FirebaseID,
		Name:                 player.Name,
		ProfileImageURL:      player.Image,
		MoxfieldUsername:     player.MoxfieldUsername,
		WinrateAllTime:       0,
		NumberofGamesAllTime: 0,
		DecksAllTime:         []DeckWithCount{},
		OpponentsAllTime:     []PlayerOpponentWithCount{},
		CurrentGame:          nil,
	}
}
func convertDeckToDto(deck *repository.Deck) Deck {
	return Deck{
		ID:           &deck.ID,
		Commander:    deck.Commander,
		Colors:       deck.Colors,
		Crop:         deck.Crop,
		SecondaryImg: deck.SecondaryImage,
		Image:        deck.Image,
		MoxfieldURL:  deck.MoxfieldURL,
		Bracket:      deck.Bracket,
		Themes:       deck.Themes,
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
