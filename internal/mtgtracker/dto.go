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
	EventType            string `json:"event_type"`
	DamageDelta          int    `json:"damage_delta"`
	TargetLifeTotalAfter int    `json:"life_total_after"`
	SourceRankingId      uint   `json:"source_ranking_id"`
	TargetRankingId      uint   `json:"target_ranking_id"`
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
}

type Ranking struct {
	PlayerID       uint   `json:"player_id"`
	Position       int    `json:"position"`
	Commander      string `json:"commander"`
	CouldHaveWon   bool   `json:"could_have_won"`
	EarlySolRing   bool   `json:"early_sol_ring"`
	StartingPlayer bool   `json:"starting_player"`
	Deck           Deck   `json:"deck"`
}

type Deck struct {
	ID           uint   `json:"id"`
	Commander    string `json:"commander"`
	Crop         string `json:"crop"`
	SecondaryImg string `json:"secondary_image"`
	Image        string `json:"image"`
}
