package opponents

import (
	"errors"
	"mtgtracker/internal/core"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	err := db.AutoMigrate(&Opponent{})
	if err != nil {
		panic("Failed to migrate Opponent repo: " + err.Error())
	}
	return &Repository{DB: db}
}

type Opponent struct {
	gorm.Model
	Player1ID  string `gorm:"not null" json:"player1_id"`
	Player2ID  string `gorm:"not null" json:"player2_id"`
	GameCount  int    `gorm:"default:0;not null" json:"game_count"` // Number of games played together

	Player1 core.Player `gorm:"foreignKey:Player1ID;references:FirebaseID" json:"player1"`
	Player2 core.Player `gorm:"foreignKey:Player2ID;references:FirebaseID" json:"player2"`
}

func (r *Repository) CreateOpponent(player1ID, player2ID string) (*Opponent, error) {
	if player1ID == player2ID {
		return nil, errors.New("cannot add yourself as opponent")
	}

	// Ensure consistent ordering for symmetrical relationship
	if player1ID > player2ID {
		player1ID, player2ID = player2ID, player1ID
	}

	opponent := Opponent{
		Player1ID: player1ID,
		Player2ID: player2ID,
	}

	if err := r.DB.Create(&opponent).Error; err != nil {
		return nil, err
	}

	// Load the related players
	if err := r.DB.Preload("Player1").Preload("Player2").First(&opponent, opponent.ID).Error; err != nil {
		return nil, err
	}

	return &opponent, nil
}

func (r *Repository) DeleteOpponent(player1ID, player2ID string) error {
	if player1ID == player2ID {
		return errors.New("invalid opponent relationship")
	}

	// Ensure consistent ordering for symmetrical relationship
	if player1ID > player2ID {
		player1ID, player2ID = player2ID, player1ID
	}

	result := r.DB.Where("player1_id = ? AND player2_id = ?", player1ID, player2ID).Delete(&Opponent{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("opponent relationship not found")
	}

	return nil
}

// OpponentWithCount represents a opponent relationship with the game count
type OpponentWithCount struct {
	Player    core.Player
	GameCount int
}

func (r *Repository) GetOpponents(playerID string, limit, offset int) ([]OpponentWithCount, int64, error) {
	var opponents []Opponent
	var total int64

	// Count total opponents
	query := r.DB.Model(&Opponent{}).Where("player1_id = ? OR player2_id = ?", playerID, playerID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated opponents where the player is either player1 or player2
	err := r.DB.Preload("Player1").Preload("Player2").
		Where("player1_id = ? OR player2_id = ?", playerID, playerID).
		Order("game_count DESC").
		Limit(limit).
		Offset(offset).
		Find(&opponents).Error
	if err != nil {
		return nil, 0, err
	}

	var opponentedPlayers []OpponentWithCount
	for _, opponent := range opponents {
		if opponent.Player1ID == playerID {
			opponentedPlayers = append(opponentedPlayers, OpponentWithCount{
				Player:    opponent.Player2,
				GameCount: opponent.GameCount,
			})
		} else {
			opponentedPlayers = append(opponentedPlayers, OpponentWithCount{
				Player:    opponent.Player1,
				GameCount: opponent.GameCount,
			})
		}
	}

	return opponentedPlayers, total, nil
}

func (r *Repository) IsOpponenting(player1ID, player2ID string) (bool, error) {
	if player1ID == player2ID {
		return false, nil
	}

	// Ensure consistent ordering for symmetrical relationship
	if player1ID > player2ID {
		player1ID, player2ID = player2ID, player1ID
	}

	var count int64
	err := r.DB.Model(&Opponent{}).Where("player1_id = ? AND player2_id = ?", player1ID, player2ID).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// IncrementGameCount increments the game count for a opponent relationship
// Creates the opponent if it doesn't exist
func (r *Repository) IncrementGameCount(player1ID, player2ID string) error {
	if player1ID == player2ID {
		return nil // Players don't opponent themselves
	}

	// Ensure consistent ordering for symmetrical relationship
	if player1ID > player2ID {
		player1ID, player2ID = player2ID, player1ID
	}

	// Try to find existing opponent
	var opponent Opponent
	err := r.DB.Where("player1_id = ? AND player2_id = ?", player1ID, player2ID).First(&opponent).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new opponent with count = 1
		opponent = Opponent{
			Player1ID: player1ID,
			Player2ID: player2ID,
			GameCount: 1,
		}
		return r.DB.Create(&opponent).Error
	} else if err != nil {
		return err
	}

	// Increment existing opponent count
	return r.DB.Model(&opponent).UpdateColumn("game_count", gorm.Expr("game_count + ?", 1)).Error
}

// DecrementGameCount decrements the game count for a opponent relationship
// Does NOT delete the opponent even if count reaches 0
func (r *Repository) DecrementGameCount(player1ID, player2ID string) error {
	if player1ID == player2ID {
		return nil // Players don't opponent themselves
	}

	// Ensure consistent ordering for symmetrical relationship
	if player1ID > player2ID {
		player1ID, player2ID = player2ID, player1ID
	}

	return r.DB.Model(&Opponent{}).
		Where("player1_id = ? AND player2_id = ?", player1ID, player2ID).
		UpdateColumn("game_count", gorm.Expr("CASE WHEN game_count > 0 THEN game_count - 1 ELSE 0 END")).
		Error
}
