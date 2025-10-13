package push

import (
	"log"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	err := db.AutoMigrate(&DeviceToken{})
	if err != nil {
		log.Fatal("Failed to migrate DeviceToken:", err)
	}
	return &Repository{DB: db}
}

// SaveToken saves or updates a device token for a player
func (r *Repository) SaveToken(playerID, token, platform string) error {
	deviceToken := DeviceToken{
		PlayerID: playerID,
		Token:    token,
		Platform: platform,
	}
	// Upsert: update if exists, create if not
	return r.DB.Where("token = ?", token).
		Assign(deviceToken).
		FirstOrCreate(&deviceToken).Error
}

// GetPlayerTokens returns all device token strings for a player (for FCM)
func (r *Repository) GetPlayerTokens(playerID string) ([]string, error) {
	var tokens []DeviceToken
	err := r.DB.Where("player_id = ?", playerID).Find(&tokens).Error
	if err != nil {
		return nil, err
	}

	result := make([]string, len(tokens))
	for i, t := range tokens {
		result[i] = t.Token
	}
	return result, nil
}

// GetPlayerDeviceTokens returns all device tokens with metadata for a player
func (r *Repository) GetPlayerDeviceTokens(playerID string) ([]DeviceToken, error) {
	var tokens []DeviceToken
	err := r.DB.Where("player_id = ?", playerID).Order("created_at DESC").Find(&tokens).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// DeleteToken removes a device token (used when token is invalid)
func (r *Repository) DeleteToken(token string) error {
	return r.DB.Where("token = ?", token).Delete(&DeviceToken{}).Error
}

// DeletePlayerToken removes a specific token for a player
func (r *Repository) DeletePlayerToken(playerID, token string) error {
	return r.DB.Where("player_id = ? AND token = ?", playerID, token).Delete(&DeviceToken{}).Error
}
