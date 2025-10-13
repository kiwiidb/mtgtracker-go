package core

import (
	"path/filepath"
	"sort"
	"time"
)

func (svc *Service) ConvertGameToDto(game *Game, addEvents bool) GameResponse {
	result := GameResponse{
		ID:         game.ID,
		CreatorID:  game.CreatorID,
		Duration:   game.Duration,
		Date:       game.Date,
		Comments:   game.Comments,
		Finished:   game.Finished,
		Rankings:   convertRankingsWithLifeTotal(game.Rankings, game.GameEvents),
		GameEvents: make([]GameEventResponse, len(game.GameEvents)),
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

func convertSimpleDeck(deck Deck) SimpleDeck {
	return SimpleDeck{
		Commander:      deck.Commander,
		Image:          deck.Image,
		SecondaryImage: deck.SecondaryImage,
		Crop:           deck.Crop,
	}
}

func convertGameEvent(event *GameEvent, uploadUrl string) GameEventResponse {
	var sourceRanking, targetRanking *RankingResponse

	if event.SourceRanking != nil {
		ranking := convertRankingToDto(event.SourceRanking)
		sourceRanking = &ranking
	}

	if event.TargetRanking != nil {
		ranking := convertRankingToDto(event.TargetRanking)
		targetRanking = &ranking
	}

	return GameEventResponse{
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

func convertDeckFromRanking(rank *Ranking) DeckResponse {
	if rank.Deck != nil {
		return DeckResponse{
			ID:           &rank.Deck.ID,
			Commander:    rank.Deck.Commander,
			Colors:       rank.Deck.Colors,
			Crop:         rank.Deck.Crop,
			SecondaryImg: rank.Deck.SecondaryImage,
			Image:        rank.Deck.Image,
			MoxfieldURL:  rank.Deck.MoxfieldURL,
			Bracket:      rank.Deck.Bracket,
		}
	}
	return DeckResponse{
		Commander:    rank.DeckEmbedded.Commander,
		Colors:       nil,
		Crop:         rank.DeckEmbedded.Crop,
		SecondaryImg: rank.DeckEmbedded.SecondaryImage,
		Image:        rank.DeckEmbedded.Image,
	}
}

func convertRankingToDto(rank *Ranking) RankingResponse {
	return RankingResponse{
		ID:       rank.ID,
		PlayerID: rank.PlayerID,
		Position: rank.Position,
		Deck:     convertDeckFromRanking(rank),
		Player: func() *PlayerResponse {
			if rank.Player != nil {
				return &PlayerResponse{
					ID:              rank.Player.FirebaseID,
					Name:            rank.Player.Name,
					ProfileImageURL: rank.Player.Image,
				}
			}
			return nil
		}(),
	}
}

func convertRankingsWithLifeTotal(rankings []Ranking, gameEvents []GameEvent) []RankingResponse {
	result := make([]RankingResponse, len(rankings))

	// Build a map of ranking ID to most recent life total event
	lifeTotalMap := make(map[uint]*GameEvent)
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

		result[i] = RankingResponse{
			ID:                     rank.ID,
			PlayerID:               rank.PlayerID,
			Position:               rank.Position,
			Deck:                   convertDeckFromRanking(&rank),
			LastLifeTotal:          lastLifeTotal,
			LastLifeTotalTimestamp: lastLifeTotalTimestamp,
			Player: func() *PlayerResponse {
				if rank.Player != nil {
					return &PlayerResponse{
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

func (svc *Service) ConvertPlayerToResponse(player *Player) PlayerResponse {
	result := PlayerResponse{
		ID:               player.FirebaseID,
		Name:             player.Name,
		ProfileImageURL:  player.Image,
		MoxfieldUsername: player.MoxfieldUsername,
	}
	// Convert player.Decks to DeckWithCount
	wins := 0
	totalGames := 0
	decks := make([]DeckWithCount, 0, len(player.Decks))
	for _, deck := range player.Decks {
		decks = append(decks, DeckWithCount{
			Deck: DeckResponse{
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
		wins += deck.WinCount
		totalGames += deck.GameCount
	}
	// sort decks by count descending
	sort.Slice(decks, func(i, j int) bool {
		return decks[i].Count > decks[j].Count
	})

	// Calculate top 2 most played colors from decks
	colors := calculateTopColors(decks, 2)
	// Calculate winrate

	var winrate float64
	if totalGames > 0 {
		winrate = float64(wins) / float64(totalGames) * 100
	}
	result.Colors = colors
	result.WinrateAllTime = winrate
	result.NumberofGamesAllTime = totalGames
	result.DecksAllTime = decks

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
func ConvertPlayerToDtoSimple(player *Player) PlayerResponse {
	return PlayerResponse{
		ID:               player.FirebaseID,
		Name:             player.Name,
		ProfileImageURL:  player.Image,
		MoxfieldUsername: player.MoxfieldUsername,
	}
}
func convertDeckToDto(deck *Deck) DeckResponse {
	return DeckResponse{
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
