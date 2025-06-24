package mtgtracker

import (
	"time"
)

type SignupPlayerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// duration int, comments, image string, rankings []Ranking
type CreateGameRequest struct {
	Duration *int       `json:"duration"`
	Date     *time.Time `json:"date"`
	Comments string     `json:"comments"`
	Image    string     `json:"image"`
	Finished bool       `json:"finished"`
	Rankings []Ranking  `json:"rankings"`
}

type UpdateGameRequest struct {
	GameID   uint      `json:"game_id"`
	Finished *bool     `json:"finished"`
	Rankings []Ranking `json:"rankings"`
}

type GameEventRequest struct {
	EventType            string  `json:"event_type"`
	EventImageName       *string `json:"event_image_name,omitempty"` // Optional field for image upload
	Comment              *string `json:"comment,omitempty"`          // Optional field for image upload
	DamageDelta          int     `json:"damage_delta"`
	TargetLifeTotalAfter int     `json:"life_total_after"`
	SourceRankingId      *uint   `json:"source_ranking_id,omitempty"` // Made nullable with pointer
	TargetRankingId      *uint   `json:"target_ranking_id,omitempty"` // Made nullable with pointer
}

type Player struct {
	ID                   uint     `json:"ID"`
	Name                 string   `json:"name"`
	WinrateAllTime       float64  `json:"winrate_all_time"`
	NumberofGamesAllTime int      `json:"number_of_games_all_time"`
	DecksAllTime         []Deck   `json:"decks_all_time"`
	CoPlayersAllTime     []Player `json:"co_players_all_time"`
	Games                []Game   `json:"games"`
	CurrentGame          *Game    `json:"current_game,omitempty"`
}

type Game struct {
	ID         uint
	Duration   *int
	Date       *time.Time
	Comments   string
	Image      string
	Rankings   []Ranking
	Finished   bool
	GameEvents []GameEvent
}

type GameEvent struct {
	GameID               uint      `json:"game_id"`
	EventType            string    `json:"event_type"`
	DamageDelta          int       `json:"damage_delta"`
	CreatedAt            time.Time `json:"created_at"`
	TargetLifeTotalAfter int       `json:"target_life_total_after"`
	SourcePlayer         string    `json:"source_player"`
	TargetPlayer         string    `json:"target_player"`
	SourceCommander      string    `json:"source_commander"`
	TargetCommander      string    `json:"target_commander"`
	ImageUrl             string    `json:"image_url"`                  // URL of the uploaded image
	UploadImageUrl       string    `json:"upload_image_url,omitempty"` // Presigned URL for image upload
}

type Ranking struct {
	ID        uint   `json:"ID"`
	PlayerID  uint   `json:"player_id"`
	Position  int    `json:"position"`
	LifeTotal *uint  `json:"life_total,omitempty"`
	Deck      Deck   `json:"deck"`
	Player    Player `json:"player,omitempty"` // Optional, can be omitted if not needed
}

type Deck struct {
	ID           uint   `json:"ID"`
	Commander    string `json:"commander"`
	Crop         string `json:"crop"`
	SecondaryImg string `json:"secondary_image"`
	Image        string `json:"image"`
}
