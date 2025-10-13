package push

import "time"

type DeviceTokenResponse struct {
	ID        uint      `json:"id"`
	Token     string    `json:"token"`
	Platform  string    `json:"platform"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toDeviceTokenResponse(token *DeviceToken) DeviceTokenResponse {
	return DeviceTokenResponse{
		ID:        token.ID,
		Token:     token.Token,
		Platform:  token.Platform,
		CreatedAt: token.CreatedAt,
		UpdatedAt: token.UpdatedAt,
	}
}
