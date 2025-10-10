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
	Player1ID string `gorm:"not null" json:"player1_id"`
	Player2ID string `gorm:"not null" json:"player2_id"`

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

func (r *Repository) GetFollows(playerID string) ([]core.Player, error) {
	var follows []Follow

	// Get all follows where the player is either player1 or player2
	err := r.DB.Preload("Player1").Preload("Player2").
		Where("player1_id = ? OR player2_id = ?", playerID, playerID).
		Find(&follows).Error
	if err != nil {
		return nil, err
	}

	var followedPlayers []core.Player
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
