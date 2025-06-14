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
		ImageUrl:               uploadUrl,
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
