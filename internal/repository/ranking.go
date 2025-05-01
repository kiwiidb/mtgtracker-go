package repository

import "sort"

func (r *Repository) GetGroupRankingByWins(groupID uint) ([]PlayerWin, error) {
	var games []Game
	err := r.DB.Preload("Rankings").Where("group_id = ?", groupID).Find(&games).Error
	if err != nil {
		return nil, err
	}

	// Maps: playerID -> total wins, and playerID -> deckID -> deckWins
	winCounts := make(map[uint]int)
	deckBreakdown := make(map[uint]map[uint]int) // map[playerID][deckID] = wins

	for _, game := range games {
		for _, rank := range game.Rankings {
			if rank.Position == 1 {
				winCounts[rank.PlayerID]++

				if _, ok := deckBreakdown[rank.PlayerID]; !ok {
					deckBreakdown[rank.PlayerID] = make(map[uint]int)
				}
				deckBreakdown[rank.PlayerID][rank.DeckID]++
				break
			}
		}
	}

	var result []PlayerWin
	for playerID, totalWins := range winCounts {
		var player Player
		if err := r.DB.First(&player, playerID).Error; err != nil {
			return nil, err
		}

		var deckWins []DeckWin
		for deckID, count := range deckBreakdown[playerID] {
			var deck Deck
			if err := r.DB.First(&deck, deckID).Error; err != nil {
				return nil, err
			}
			deckWins = append(deckWins, DeckWin{
				DeckID: deckID,
				Deck:   deck,
				Wins:   count,
			})
		}

		result = append(result, PlayerWin{
			PlayerID: playerID,
			Player:   player,
			Wins:     totalWins,
			DeckWins: deckWins,
		})
	}

	// Sort by total wins descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].Wins > result[j].Wins
	})

	return result, nil
}
