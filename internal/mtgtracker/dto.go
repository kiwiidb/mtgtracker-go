package mtgtracker

import (
	"time"
)

type SignupPlayerRequest struct {
	Name string `json:"name"`
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

type PlayerWithCount struct {
	Player Player `json:"player"`
	Count  int    `json:"count"`
}
type DeckWithCount struct {
	Deck  Deck `json:"deck"`
	Count int  `json:"count"`
	Wins  int  `json:"wins"`
}
type Player struct {
	ID                   uint              `json:"id"`
	Name                 string            `json:"name"`
	WinrateAllTime       float64           `json:"winrate_all_time"`
	NumberofGamesAllTime int               `json:"number_of_games_all_time"`
	DecksAllTime         []DeckWithCount   `json:"decks_all_time"`
	CoPlayersAllTime     []PlayerWithCount `json:"co_players_all_time"`
	Games                []Game            `json:"games"`
	CurrentGame          *Game             `json:"current_game,omitempty"`
}

type Game struct {
	ID         uint        `json:"id"`
	Duration   *int        `json:"duration,omitempty"`
	Date       *time.Time  `json:"date,omitempty"`
	Comments   string      `json:"comments,omitempty"`
	Rankings   []Ranking   `json:"rankings,omitempty"`
	Finished   bool        `json:"finished"`
	GameEvents []GameEvent `json:"game_events,omitempty"`
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
	ID        uint          `json:"id"`
	PlayerID  uint          `json:"player_id"`
	Position  int           `json:"position"`
	LifeTotal *uint         `json:"life_total,omitempty"`
	Deck      Deck          `json:"deck"`
	Player    Player        `json:"player,omitempty"` // Optional, can be omitted if not needed
	Status    RankingStatus `json:"status,omitempty"` // Optional, can be omitted if not needed
}

type RankingStatus string

const (
	StatusPending  RankingStatus = "pending"
	StatusAccepted RankingStatus = "accepted"
	StatusDeclined RankingStatus = "declined"
)

type Deck struct {
	Commander    string `json:"commander"`
	Crop         string `json:"crop"`
	SecondaryImg string `json:"secondary_image"`
	Image        string `json:"image"`
}
