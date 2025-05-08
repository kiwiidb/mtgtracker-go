package mtgtracker

import "time"

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
	Rankings []Ranking  `json:"rankings"`
}

type Ranking struct {
	PlayerID     uint   `json:"player_id"`
	Commander    string `json:"commander"`
	Position     int    `json:"position"`
	CouldHaveWon bool   `json:"could_have_won"`
}
