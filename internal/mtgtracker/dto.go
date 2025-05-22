package mtgtracker

import (
	"time"
)

type SignupPlayerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Image string `json:"image"`
}

type CreateGroupRequest struct {
	CreatorID uint   `json:"creator_id"`
	Name      string `json:"name"`
	Image     string `json:"image"`
}

type AddDeckToPlayerRequest struct {
	MoxfieldURL string `json:"moxfield_url"`
	PlayerID    uint   `json:"player_id"`
	Commander   string `json:"commander"`
}

// groupID uint, duration int, comments, image string, rankings []Ranking
type CreateGameRequest struct {
	GroupID  uint       `json:"group_id"`
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
	Id         uint
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
	Commander      string `json:"commander"`
	CommanderImage string `json:"commander_image"`
	Position       int    `json:"position"`
	CouldHaveWon   bool   `json:"could_have_won"`
	EarlySolRing   bool   `json:"early_sol_ring"`
	StartingPlayer bool   `json:"starting_player"`
}
