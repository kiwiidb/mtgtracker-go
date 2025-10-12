package follows

import (
	"errors"
	"mtgtracker/internal/core"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	err := db.AutoMigrate(&Follow{})
	if err != nil {
		panic("Failed to migrate Follow repo: " + err.Error())
	}
	return &Repository{DB: db}
}

type Follow struct {
	gorm.Model
	Player1ID  string `gorm:"not null" json:"player1_id"`
	Player2ID  string `gorm:"not null" json:"player2_id"`
	GameCount  int    `gorm:"default:0;not null" json:"game_count"` // Number of games played together

	Player1 core.Player `gorm:"foreignKey:Player1ID;references:FirebaseID" json:"player1"`
	Player2 core.Player `gorm:"foreignKey:Player2ID;references:FirebaseID" json:"player2"`
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

// FollowWithCount represents a follow relationship with the game count
type FollowWithCount struct {
	Player    core.Player
	GameCount int
}

func (r *Repository) GetFollows(playerID string) ([]FollowWithCount, error) {
	var follows []Follow

	// Get all follows where the player is either player1 or player2
	err := r.DB.Preload("Player1").Preload("Player2").
		Where("player1_id = ? OR player2_id = ?", playerID, playerID).
		Find(&follows).Error
	if err != nil {
		return nil, err
	}

	var followedPlayers []FollowWithCount
	for _, follow := range follows {
		if follow.Player1ID == playerID {
			followedPlayers = append(followedPlayers, FollowWithCount{
				Player:    follow.Player2,
				GameCount: follow.GameCount,
			})
		} else {
			followedPlayers = append(followedPlayers, FollowWithCount{
				Player:    follow.Player1,
				GameCount: follow.GameCount,
			})
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

// IncrementGameCount increments the game count for a follow relationship
// Creates the follow if it doesn't exist
func (r *Repository) IncrementGameCount(player1ID, player2ID string) error {
	if player1ID == player2ID {
		return nil // Players don't follow themselves
	}

	// Ensure consistent ordering for symmetrical relationship
	if player1ID > player2ID {
		player1ID, player2ID = player2ID, player1ID
	}

	// Try to find existing follow
	var follow Follow
	err := r.DB.Where("player1_id = ? AND player2_id = ?", player1ID, player2ID).First(&follow).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new follow with count = 1
		follow = Follow{
			Player1ID: player1ID,
			Player2ID: player2ID,
			GameCount: 1,
		}
		return r.DB.Create(&follow).Error
	} else if err != nil {
		return err
	}

	// Increment existing follow count
	return r.DB.Model(&follow).UpdateColumn("game_count", gorm.Expr("game_count + ?", 1)).Error
}

// DecrementGameCount decrements the game count for a follow relationship
// Does NOT delete the follow even if count reaches 0
func (r *Repository) DecrementGameCount(player1ID, player2ID string) error {
	if player1ID == player2ID {
		return nil // Players don't follow themselves
	}

	// Ensure consistent ordering for symmetrical relationship
	if player1ID > player2ID {
		player1ID, player2ID = player2ID, player1ID
	}

	return r.DB.Model(&Follow{}).
		Where("player1_id = ? AND player2_id = ?", player1ID, player2ID).
		UpdateColumn("game_count", gorm.Expr("CASE WHEN game_count > 0 THEN game_count - 1 ELSE 0 END")).
		Error
}
