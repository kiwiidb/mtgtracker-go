package mtgtracker

import (
	"time"
)

type SignupPlayerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Image string `json:"image"`
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
	DamageDelta          int     `json:"damage_delta"`
	TargetLifeTotalAfter int     `json:"life_total_after"`
	SourceRankingId      *uint   `json:"source_ranking_id,omitempty"` // Made nullable with pointer
	TargetRankingId      *uint   `json:"target_ranking_id,omitempty"` // Made nullable with pointer
}

type PlayerDto struct {
	ID   uint   `json:"ID"`
	Name string `json:"name"`
	//WinrateMonthly    float64 `json:"winrate_monthly"`
	WinrateAllTime       float64 `json:"winrate_all_time"`
	NumberofGamesAllTime int     `json:"number_of_games_all_time"`
	//CommandersMonthly []Deck  `json:"commanders_monthly"`
	DecksAllTime     []Deck    `json:"decks_all_time"`
	CoPlayersAllTime []Player  `json:"co_players_all_time"`
	Games            []GameDto `json:"games"`
}

type GameDto struct {
	ID         uint
	Duration   *int
	Date       *time.Time
	Comments   string
	Image      string
	Rankings   []Ranking
	Finished   bool
	GameEvents []GameEventDto
}

type GameEventDto struct {
	GameID                 uint
	EventType              string
	DamageDelta            int
	CreatedAt              time.Time
	TargetLifeTotalAfter   int
	SourcePlayer           string
	TargetPlayer           string
	SourceCommanderCropImg string
	TargetCommanderCropImg string
	SourceCommander        string
	TargetCommander        string
	ImageUrl               string // New field for uploaded image URL
	UploadImageUrl         string // Optional, can be used for image upload
}

type Ranking struct {
	ID             uint   `json:"ID"`
	PlayerID       uint   `json:"player_id"`
	Position       int    `json:"position"`
	Commander      string `json:"commander"`
	CouldHaveWon   bool   `json:"could_have_won"`
	EarlySolRing   bool   `json:"early_sol_ring"`
	StartingPlayer bool   `json:"starting_player"`
	Deck           Deck   `json:"deck"`
	Player         Player `json:"player,omitempty"` // Optional, can be omitted if not needed
}

type Player struct {
	ID   uint   `json:"ID"`
	Name string `json:"name"`
}
type Deck struct {
	ID           uint   `json:"ID"`
	Commander    string `json:"commander"`
	Crop         string `json:"crop"`
	SecondaryImg string `json:"secondary_image"`
	Image        string `json:"image"`
}
