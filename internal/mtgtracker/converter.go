package mtgtracker

import "mtgtracker/internal/repository"

func convertGameToDto(game *repository.Game) GameDto {
	result := GameDto{
		ID:         game.ID,
		Duration:   game.Duration,
		Date:       game.Date,
		Comments:   game.Comments,
		Image:      game.Image,
		Rankings:   convertRankingsToDto(game.Rankings),
		GameEvents: make([]GameEventDto, len(game.GameEvents)),
	}
	for i, event := range game.GameEvents {
		result.GameEvents[i] = GameEventDto{
			GameID:                 event.GameID,
			EventType:              event.EventType,
			DamageDelta:            event.DamageDelta,
			CreatedAt:              event.CreatedAt,
			TargetLifeTotalAfter:   event.TargetLifeTotalAfter,
			SourcePlayer:           event.SourceRanking.Player.Name,
			TargetPlayer:           event.TargetRanking.Player.Name,
			SourceCommander:        event.SourceRanking.Deck.Commander,
			TargetCommander:        event.TargetRanking.Deck.Commander,
			SourceCommanderCropImg: event.SourceRanking.Deck.Crop,
			TargetCommanderCropImg: event.TargetRanking.Deck.Crop,
		}
	}
	return result
}

func convertRankingsToDto(rankings []repository.Ranking) []Ranking {
	result := make([]Ranking, len(rankings))
	for i, rank := range rankings {
		result[i] = Ranking{
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
