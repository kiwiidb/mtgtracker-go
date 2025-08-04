package repository

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) GetPendingRankings(firebaseID string) ([]Ranking, error) {
	var rankings []Ranking
	err := r.DB.Where("status = ? AND firebase_id = ?", StatusPending, firebaseID).Preload("Player").Find(&rankings).Error
	if err != nil {
		return nil, err
	}
	return rankings, nil
}

// todo: check if ranking is for this user
func (r *Repository) AcceptRanking(rankingID uint) (*Ranking, error) {
	var ranking Ranking
	err := r.DB.First(&ranking, rankingID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ranking not found")
		}
		return nil, err
	}

	// Update the ranking status
	ranking.Status = StatusAccepted
	err = r.DB.Save(&ranking).Error
	if err != nil {
		return nil, err
	}
	
	// Get all other accepted players in the same game
	acceptedPlayerIDs, err := r.GetAcceptedPlayersInGame(ranking.GameID)
	if err != nil {
		log.Printf("Error getting accepted players for game %d: %v", ranking.GameID, err)
		// Don't fail the ranking acceptance if follow creation fails
	} else {
		// Create follows between the current player and all other accepted players
		for _, otherPlayerID := range acceptedPlayerIDs {
			if otherPlayerID != ranking.PlayerID {
				_, err := r.CreateFollow(ranking.PlayerID, otherPlayerID)
				if err != nil {
					// Log error but don't fail the ranking acceptance
					log.Printf("Error creating follow between players %d and %d: %v", ranking.PlayerID, otherPlayerID, err)
				}
			}
		}
	}
	
	return &ranking, nil
}

// todo: check if ranking is for this user
func (r *Repository) DeclineRanking(rankingID uint) error {
	var ranking Ranking
	err := r.DB.First(&ranking, rankingID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ranking not found")
		}
		return err
	}

	// Update the ranking status
	ranking.Status = StatusDeclined
	err = r.DB.Save(&ranking).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetPlayerByFirebaseID(userID string) (*Player, error) {
	var player Player
	// add games that are not finished
	err := r.DB.Preload("Games").Where("firebase_id = ?", userID).First(&player).Error
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *Repository) GetPlayers(search string) ([]Player, error) {
	var players []Player
	// If search is provided, filter players by name
	if search != "" {
		err := r.DB.Where("name LIKE ?", "%"+search+"%").Find(&players).Error
		if err != nil {
			return nil, err
		}
		return players, nil
	}
	err := r.DB.Find(&players).Error
	if err != nil {
		return nil, err
	}
	return players, nil
}

func (r *Repository) GetPlayerByID(id uint) (*Player, error) {
	var player Player
	err := r.DB.First(&player, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("player not found")
		}
		return nil, err
	}

	// Get games through rankings
	var rankings []Ranking
	err = r.DB.Preload("Player").Where("player_id = ?", id).Find(&rankings).Error
	if err != nil {
		return nil, err
	}

	// Get unique game IDs
	gameIDs := make(map[uint]bool)
	for _, ranking := range rankings {
		gameIDs[ranking.GameID] = true
	}

	// Convert map keys to slice
	gameIDSlice := make([]uint, 0, len(gameIDs))
	for gameID := range gameIDs {
		gameIDSlice = append(gameIDSlice, gameID)
	}

	// Get games with all their rankings and related data
	var games []Game
	if len(gameIDSlice) > 0 {
		err = r.DB.Preload("Rankings.Player").Preload("GameEvents").Where("id IN ?", gameIDSlice).Order("date desc").Find(&games).Error
		if err != nil {
			return nil, err
		}
	}

	player.Games = games
	return &player, nil
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
	return res, nil
}

func (r *Repository) GetGames() ([]Game, error) {
	var games []Game
	err := r.DB.Preload("Rankings", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Player")
	}).Preload("GameEvents").Order("Date desc").Find(&games).Error
	if err != nil {
		return nil, err
	}

	return games, nil
}

func (r *Repository) DeleteDeck(deckID uint) error {
	var deck Deck
	if err := r.DB.First(&deck, deckID).Error; err != nil {
		return errors.New("deck not found")
	}
	return r.DB.Delete(&deck).Error
}

func NewRepository(db *gorm.DB) *Repository {
	err := db.AutoMigrate(&Player{}, &Game{}, &Ranking{}, &GameEvent{}, &Follow{})
	if err != nil {
		log.Fatal(err)
	}
	
	// Add unique constraint for symmetrical follows
	db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_follow_pair 
		ON follows (LEAST(player1_id, player2_id), GREATEST(player1_id, player2_id))`)
	
	return &Repository{
		DB: db,
	}
}
func (r *Repository) InsertPlayer(name string, email string, userId string) (*Player, error) {
	player := Player{Name: name, Email: email, FirebaseID: userId}
	result := r.DB.Create(&player)
	return &player, result.Error
}

func (r *Repository) InsertGame(duration *int, comments, image string, date *time.Time, finished bool, rankings []Ranking) (*Game, error) {

	// Ensure each ranking has valid player and deck (optional but safe)
	for i, rank := range rankings {
		var player Player
		if err := r.DB.First(&player, rank.PlayerID).Error; err != nil {
			return nil, errors.New("invalid player ID")
		}
		rank.Player = player
		rankings[i] = rank

	}

	game := Game{
		Duration: duration,
		Date:     date,
		Comments: comments,
		Image:    image,
		Rankings: rankings,
	}

	if err := r.DB.Create(&game).Error; err != nil {
		return nil, err
	}

	return &game, nil
}

func (r *Repository) InsertGameEvent(gameId uint, eventType string, damageDelta, lifeAfter int, source, target *uint, imageUrl string) (*GameEvent, error) {
	event := GameEvent{
		GameID:               gameId,
		EventType:            eventType,
		DamageDelta:          damageDelta,
		TargetLifeTotalAfter: lifeAfter,
		SourceRankingID:      source,
		TargetRankingID:      target,
		ImageUrl:             imageUrl,
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
		Preload("GameEvents.SourceRanking.Player").
		Preload("GameEvents.TargetRanking.Player").
		First(&game, gameID).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *Repository) CreateFollow(player1ID, player2ID uint) (*Follow, error) {
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

func (r *Repository) DeleteFollow(player1ID, player2ID uint) error {
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

func (r *Repository) GetFollows(playerID uint) ([]Player, error) {
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

func (r *Repository) IsFollowing(player1ID, player2ID uint) (bool, error) {
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

func (r *Repository) GetAcceptedPlayersInGame(gameID uint) ([]uint, error) {
	var rankings []Ranking
	err := r.DB.Where("game_id = ? AND status = ?", gameID, StatusAccepted).Find(&rankings).Error
	if err != nil {
		return nil, err
	}
	
	playerIDs := make([]uint, 0, len(rankings))
	for _, ranking := range rankings {
		playerIDs = append(playerIDs, ranking.PlayerID)
	}
	
	return playerIDs, nil
}
