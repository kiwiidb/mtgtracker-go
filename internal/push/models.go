package push

import (
	"mtgtracker/internal/core"

	"gorm.io/gorm"
)

type DeviceToken struct {
	gorm.Model
	PlayerID string `gorm:"not null" json:"player_id"`
	Token    string `gorm:"not null;unique" json:"token"`
	Platform string `gorm:"not null" json:"platform"` // "ios", "android", "web"

	Player *core.Player `gorm:"foreignKey:PlayerID;references:FirebaseID" json:"player,omitempty"`
}
