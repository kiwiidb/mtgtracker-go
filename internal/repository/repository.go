package repository

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	DB       *gorm.DB
	notifier Notifier
}

type Notifier interface {
	gameCreated(game *Game) error
	gameFinished(game *Game) error
}

func (r *Repository) GetPlayerByFirebaseID(userID string) (*Player, error) {
	var player Player
	err := r.DB.Preload("Decks").Where("firebase_id = ?", userID).First(&player).Error
	if err != nil {
		return nil, err
	}

	// Manually load games for this player since Games field has gorm:"-"
	// We need to find games where this player has rankings
	var games []Game
	err = r.DB.Joins("JOIN rankings ON games.id = rankings.game_id").
		Where("rankings.player_id = ?", player.FirebaseID).
		Preload("Rankings", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Player").Preload("Deck")
		}).
		Preload("GameEvents").
		Distinct().
		Order("games.created_at desc").
		Find(&games).Error
	if err != nil {
		return nil, err
	}

	// Populate the Games field manually
	player.Games = games

	return &player, nil
}

func (r *Repository) GetPlayers(search string, limit, offset int) ([]Player, int64, error) {
	var players []Player
	var total int64

	query := r.DB.Model(&Player{})

	// If search is provided, filter players by name (case insensitive)
	if search != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+search+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("name ASC").Limit(limit).Offset(offset).Find(&players).Error
	if err != nil {
		return nil, 0, err
	}

	return players, total, nil
}

func (r *Repository) GetActiveGameForPlayer(playerID string) (*Game, error) {
	var game Game
	err := r.DB.
		Joins("JOIN rankings ON games.id = rankings.game_id").
		Where("rankings.player_id = ?", playerID).
		Where("games.finished = ?", false).
		Preload("Rankings", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Player").Preload("Deck")
		}).
		Preload("GameEvents.SourceRanking.Player").
		Preload("GameEvents.TargetRanking.Player").
		Preload("Creator").
		Order("games.created_at DESC").
		First(&game).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No active game found, return nil without error
		}
		return nil, err
	}
	return &game, nil
}

func (r *Repository) DeleteGame(gameID uint) error {
	// Use a transaction to ensure all deletions succeed or fail together
	return r.DB.Transaction(func(tx *gorm.DB) error {
		// First delete game events
		if err := tx.Where("game_id = ?", gameID).Delete(&GameEvent{}).Error; err != nil {
			return err
		}

		// Then delete rankings
		if err := tx.Where("game_id = ?", gameID).Delete(&Ranking{}).Error; err != nil {
			return err
		}

		// Finally delete the game itself
		if err := tx.Delete(&Game{}, gameID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *Repository) UpdateGame(gameId uint, rankings []Ranking, finished *bool) (*Game, error) {
	// update the rankings

	for _, rank := range rankings {
		log.Println("updating ranking ", rank.ID, rank.PlayerID, rank.Position)
		if err := r.DB.Model(&rank).Where("id = ?", rank.ID).Updates(rank).Error; err != nil {
			return nil, err
		}
	}
	// update the game as finished
	if finished != nil {
		if err := r.DB.Model(&Game{}).Where("id = ?", gameId).Update("finished", *finished).Error; err != nil {
			return nil, err
		}
	}
	// find the game with rankings and return it
	res, err := r.GetGameWithEvents(gameId)
	if err != nil {
		return nil, err
	}

	// Create finished game notifications if game was just finished
	if finished != nil && *finished {
		// Update deck statistics when game is finished
		if err := r.updateDeckStatisticsOnFinish(res); err != nil {
			log.Printf("Failed to update deck statistics: %v", err)
			// Don't fail the game update if deck stats update fails
		}

		if err := r.notifier.gameFinished(res); err != nil {
			log.Printf("Failed to create finished game notifications: %v", err)
			// Don't fail the game update if notifications fail
		}
	}

	return res, nil
}

func (r *Repository) GetGames(limit, offset int) ([]Game, int64, error) {
	var games []Game
	var total int64

	// Get total count
	if err := r.DB.Model(&Game{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.DB.Preload("Rankings", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Player").Preload("Deck")
	}).Preload("GameEvents").Order("Date desc").Limit(limit).Offset(offset).Find(&games).Error
	if err != nil {
		return nil, 0, err
	}

	return games, total, nil
}

func (r *Repository) GetPlayerGames(playerID string, limit, offset int) ([]Game, int64, error) {
	var games []Game
	var total int64

	// Subquery to get game IDs where player participated
	subQuery := r.DB.Model(&Ranking{}).Select("game_id").Where("player_id = ?", playerID)

	// Get total count
	if err := r.DB.Model(&Game{}).Where("id IN (?)", subQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.DB.Where("id IN (?)", subQuery).
		Preload("Rankings", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Player").Preload("Deck")
		}).
		Preload("GameEvents").
		Preload("Creator").
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&games).Error
	if err != nil {
		return nil, 0, err
	}

	return games, total, nil
}

func NewRepository(db *gorm.DB) *Repository {
	err := db.AutoMigrate(&Player{}, &Game{}, &Ranking{}, &GameEvent{}, &Follow{}, &Deck{})
	if err != nil {
		log.Fatal(err)
	}

	// Add unique constraint for symmetrical follows
	db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_follow_pair
		ON follows (player1_id, player2_id)`)

	return &Repository{
		DB: db,
	}
}
func (r *Repository) InsertPlayer(name string, email string, userId string) (*Player, error) {
	player := Player{Name: name, Email: email, FirebaseID: userId}
	result := r.DB.Create(&player)
	return &player, result.Error
}

func (r *Repository) InsertGame(creatorID string, duration *int, comments, image string, date *time.Time, finished bool, rankings []Ranking) (*Game, error) {

	// Ensure each ranking has valid player and deck
	for i, rank := range rankings {
		if rank.PlayerID != nil {
			var player Player
			if err := r.DB.Where("firebase_id = ?", *rank.PlayerID).First(&player).Error; err != nil {
				return nil, errors.New("invalid player ID")
			}
			rank.Player = &player
		} else {
			rank.Player = nil
		}

		// If DeckID is provided, validate and load the deck
		if rank.DeckID != nil {
			log.Println("DECKID PROVIDED:", *rank.DeckID)
			var deck Deck
			if err := r.DB.First(&deck, *rank.DeckID).Error; err != nil {
				return nil, errors.New("invalid deck ID")
			}
			// Verify deck belongs to the player
			if rank.PlayerID != nil && deck.PlayerID != nil && *deck.PlayerID != *rank.PlayerID {
				return nil, errors.New("deck does not belong to player")
			}
			rank.Deck = &deck
		}

		rankings[i] = rank
	}

	game := Game{
		CreatorID: &creatorID,
		Duration:  duration,
		Date:      date,
		Comments:  comments,
		Image:     image,
		Rankings:  rankings,
	}

	if err := r.DB.Create(&game).Error; err != nil {
		return nil, err
	}

	// Create notifications for all players in the game
	if err := r.notifier.gameCreated(&game); err != nil {
		log.Printf("Failed to create game notifications: %v", err)
		// Don't fail the game creation if notifications fail
	}
	if err := r.createInitGameEvents(&game); err != nil {
		log.Printf("Failed to create initial game events: %v", err)
		// Don't fail the game creation if initial game events fail
	}

	return &game, nil
}

func (r *Repository) createInitGameEvents(game *Game) error {
	for _, ranking := range game.Rankings {
		event := GameEvent{
			GameID:               game.ID,
			EventType:            EventTypeInit,
			TargetRankingID:      &ranking.ID,
			TargetLifeTotalAfter: 40, // Default starting life total for Commander
		}
		if err := r.DB.Create(&event).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) InsertGameEvent(gameId uint, eventType string, damageDelta, lifeAfter int, source, target *uint, imageUrl string, comment *string) (*GameEvent, error) {
	event := GameEvent{
		GameID:               gameId,
		EventType:            eventType,
		DamageDelta:          damageDelta,
		TargetLifeTotalAfter: lifeAfter,
		SourceRankingID:      source,
		TargetRankingID:      target,
		ImageUrl:             imageUrl,
		Comment:              comment,
	}
	if err := r.DB.Create(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *Repository) GetGameWithEvents(gameID uint) (*Game, error) {
	var game Game
	err := r.DB.
		Preload("Rankings.Player").
		Preload("Rankings.Deck").
		Preload("GameEvents.SourceRanking.Player").
		Preload("GameEvents.SourceRanking.Deck").
		Preload("GameEvents.TargetRanking.Player").
		Preload("GameEvents.TargetRanking.Deck").
		First(&game, gameID).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *Repository) CreateFollow(player1ID, player2ID string) (*Follow, error) {
	if player1ID == player2ID {
		return nil, errors.New("cannot follow yourself")
	}

	// Ensure consistent ordering for symmetrical relationship
	if player1ID > player2ID {
		player1ID, player2ID = player2ID, player1ID
	}

	follow := Follow{
		Player1ID: player1ID,
		Player2ID: player2ID,
	}

	if err := r.DB.Create(&follow).Error; err != nil {
		return nil, err
	}

	// Load the related players
	if err := r.DB.Preload("Player1").Preload("Player2").First(&follow, follow.ID).Error; err != nil {
		return nil, err
	}

	return &follow, nil
}

func (r *Repository) DeleteFollow(player1ID, player2ID string) error {
	if player1ID == player2ID {
		return errors.New("invalid follow relationship")
	}

	// Ensure consistent ordering for symmetrical relationship
	if player1ID > player2ID {
		player1ID, player2ID = player2ID, player1ID
	}

	result := r.DB.Where("player1_id = ? AND player2_id = ?", player1ID, player2ID).Delete(&Follow{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("follow relationship not found")
	}

	return nil
}

func (r *Repository) GetFollows(playerID string) ([]Player, error) {
	var follows []Follow

	// Get all follows where the player is either player1 or player2
	err := r.DB.Preload("Player1").Preload("Player2").
		Where("player1_id = ? OR player2_id = ?", playerID, playerID).
		Find(&follows).Error
	if err != nil {
		return nil, err
	}

	var followedPlayers []Player
	for _, follow := range follows {
		if follow.Player1ID == playerID {
			followedPlayers = append(followedPlayers, follow.Player2)
		} else {
			followedPlayers = append(followedPlayers, follow.Player1)
		}
	}

	return followedPlayers, nil
}

func (r *Repository) IsFollowing(player1ID, player2ID string) (bool, error) {
	if player1ID == player2ID {
		return false, nil
	}

	// Ensure consistent ordering for symmetrical relationship
	if player1ID > player2ID {
		player1ID, player2ID = player2ID, player1ID
	}

	var count int64
	err := r.DB.Model(&Follow{}).Where("player1_id = ? AND player2_id = ?", player1ID, player2ID).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Repository) UpdatePlayerProfileImage(firebaseID, imageURL string) error {
	result := r.DB.Model(&Player{}).Where("firebase_id = ?", firebaseID).Update("image", imageURL)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("player not found")
	}
	return nil
}

func (r *Repository) UpdatePlayer(firebaseID string, updates map[string]interface{}) (*Player, error) {
	result := r.DB.Model(&Player{}).Where("firebase_id = ?", firebaseID).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("player not found")
	}

	// Fetch and return the updated player
	var player Player
	if err := r.DB.Where("firebase_id = ?", firebaseID).First(&player).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *Repository) DeleteRanking(rankingID uint, userID string) error {
	// First, get the ranking to verify it exists and get the game info
	var ranking Ranking
	if err := r.DB.First(&ranking, rankingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ranking not found")
		}
		return err
	}

	// Check if ranking is already deleted (player_id is nil)
	if ranking.PlayerID == nil {
		return errors.New("ranking already deleted")
	}

	// Get the game to verify the user is authorized (must be creator or the player in the ranking)
	var game Game
	if err := r.DB.First(&game, ranking.GameID).Error; err != nil {
		return err
	}

	// Check authorization: user must be either the game creator or the player in the ranking
	isCreator := game.CreatorID != nil && *game.CreatorID == userID
	isRankingPlayer := ranking.PlayerID != nil && *ranking.PlayerID == userID
	if !isCreator && !isRankingPlayer {
		return errors.New("unauthorized to delete this ranking")
	}

	// Decrement deck statistics if this ranking has a deck reference
	if ranking.DeckID != nil && game.Finished {
		if err := r.decrementDeckStatistics(*ranking.DeckID, ranking.Position == 1); err != nil {
			log.Printf("Failed to decrement deck statistics: %v", err)
			// Don't fail the deletion if deck stats update fails
		}
	}

	// Set player_id to nil instead of deleting the ranking
	return r.DB.Model(&ranking).Update("player_id", nil).Error
}

func (r *Repository) updateDeckStatisticsOnFinish(game *Game) error {
	for _, ranking := range game.Rankings {
		// Only update if the ranking has a deck reference
		if ranking.DeckID != nil {
			isWinner := ranking.Position == 1
			if err := r.incrementDeckStatistics(*ranking.DeckID, isWinner); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Repository) incrementDeckStatistics(deckID uint, isWinner bool) error {
	log.Println("Incrementing stats for deck", deckID, "winner:", isWinner)
	updates := map[string]interface{}{
		"game_count": gorm.Expr("game_count + ?", 1),
	}
	if isWinner {
		updates["win_count"] = gorm.Expr("win_count + ?", 1)
	}
	return r.DB.Model(&Deck{}).Where("id = ?", deckID).Updates(updates).Error
}

func (r *Repository) decrementDeckStatistics(deckID uint, wasWinner bool) error {
	updates := map[string]interface{}{
		"game_count": gorm.Expr("CASE WHEN game_count > 0 THEN game_count - 1 ELSE 0 END"),
	}
	if wasWinner {
		updates["win_count"] = gorm.Expr("CASE WHEN win_count > 0 THEN win_count - 1 ELSE 0 END")
	}
	return r.DB.Model(&Deck{}).Where("id = ?", deckID).Updates(updates).Error
}

func (r *Repository) CreateDeck(playerID, commander, image, secondaryImage, crop string, moxFieldID *string, themes, colors []string, bracket uint) (*Deck, error) {
	deck := Deck{
		PlayerID:       &playerID,
		MoxfieldURL:    moxFieldID,
		Themes:         themes,
		Bracket:        bracket,
		Commander:      commander,
		Colors:         colors,
		Image:          image,
		SecondaryImage: secondaryImage,
		Crop:           crop,
	}

	if err := r.DB.Create(&deck).Error; err != nil {
		return nil, err
	}

	// Load the player relationship
	if err := r.DB.Preload("Player").First(&deck, deck.ID).Error; err != nil {
		return nil, err
	}

	return &deck, nil
}

func (r *Repository) GetPlayerDecks(playerID string, limit, offset int) ([]Deck, int64, error) {
	var decks []Deck
	var total int64

	query := r.DB.Model(&Deck{}).Where("player_id = ?", playerID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.DB.Where("player_id = ?", playerID).
		Order("game_count DESC, win_count DESC").
		Limit(limit).Offset(offset).
		Find(&decks).Error
	if err != nil {
		return nil, 0, err
	}
	return decks, total, nil
}
