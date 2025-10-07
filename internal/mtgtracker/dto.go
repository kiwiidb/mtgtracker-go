package mtgtracker

import (
	"net/http"
	"strconv"
	"time"
)

type Pagination struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func (p *Pagination) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage <= 0 {
		p.PerPage = 10
	} else if p.PerPage > 100 {
		p.PerPage = 100
	}
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func ParsePagination(r *http.Request) Pagination {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	p := Pagination{Page: page, PerPage: perPage}
	p.Normalize()
	return p
}

type PaginatedResult[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
}

type SignupPlayerRequest struct {
	Name string `json:"name"`
}

// duration int, comments, image string, rankings []Ranking
type CreateGameRequest struct {
	Duration *int                   `json:"duration"`
	Date     *time.Time             `json:"date"`
	Comments string                 `json:"comments"`
	Image    string                 `json:"image"`
	Finished bool                   `json:"finished"`
	Rankings []CreateRankingRequest `json:"rankings"`
}

type CreateRankingRequest struct {
	PlayerID *string `json:"player_id,omitempty"`
	DeckID   *uint   `json:"deck_id,omitempty"` // Optional: reference to existing deck
	Deck     *Deck   `json:"deck,omitempty"`    // Optional: inline deck info (used if deck_id not provided)
}
type UpdateGameRequest struct {
	GameID   uint            `json:"game_id"`
	Finished *bool           `json:"finished"`
	Rankings []UpdateRanking `json:"rankings"`
}

type UpdateRanking struct {
	RankingID uint `json:"ranking_id"`
	Position  int  `json:"position"`
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

type PlayerOpponentWithCount struct {
	Player Player `json:"player"`
	Count  int    `json:"count"`
}
type DeckWithCount struct {
	Deck  Deck `json:"deck"`
	Count int  `json:"count"`
	Wins  int  `json:"wins"`
}
type Player struct {
	ID                   string                    `json:"id"`
	Name                 string                    `json:"name"`
	ProfileImageURL      string                    `json:"profile_image_url,omitempty"`
	MoxfieldUsername     string                    `json:"moxfield_username,omitempty"`
	Colors               []string                  `json:"colors,omitempty"` // Top 2 most played colors
	WinrateAllTime       float64                   `json:"winrate_all_time"`
	NumberofGamesAllTime int                       `json:"number_of_games_all_time"`
	DecksAllTime         []DeckWithCount           `json:"decks_all_time"`
	OpponentsAllTime     []PlayerOpponentWithCount `json:"opponents_all_time"`
	CurrentGame          *Game                     `json:"current_game,omitempty"`
}

type UpdatePlayerRequest struct {
	MoxfieldUsername *string `json:"moxfield_username,omitempty"`
}

type Game struct {
	ID         uint        `json:"id"`
	CreatorID  *string     `json:"creator_id,omitempty"`
	Duration   *int        `json:"duration,omitempty"`
	Date       *time.Time  `json:"date,omitempty"`
	Comments   string      `json:"comments,omitempty"`
	Rankings   []Ranking   `json:"rankings,omitempty"`
	Finished   bool        `json:"finished"`
	GameEvents []GameEvent `json:"game_events,omitempty"`
	Creator    *Player     `json:"creator,omitempty"`
}

type GameEvent struct {
	GameID               uint      `json:"game_id"`
	EventType            string    `json:"event_type"`
	DamageDelta          int       `json:"damage_delta"`
	CreatedAt            time.Time `json:"created_at"`
	TargetLifeTotalAfter int       `json:"target_life_total_after"`
	SourceRanking        *Ranking  `json:"source_ranking,omitempty"`
	TargetRanking        *Ranking  `json:"target_ranking,omitempty"`
	ImageUrl             string    `json:"image_url"`                  // URL of the uploaded image
	UploadImageUrl       string    `json:"upload_image_url,omitempty"` // Presigned URL for image upload
	Comment              *string   `json:"comment,omitempty"`          // New field for text description
}

type Ranking struct {
	ID                     uint       `json:"id"`
	PlayerID               *string    `json:"player_id,omitempty"`
	Position               int        `json:"position"`
	LastLifeTotal          *int       `json:"last_life_total,omitempty"`
	LastLifeTotalTimestamp *time.Time `json:"last_life_total_timestamp,omitempty"`
	Deck                   Deck       `json:"deck"`
	Player                 *Player    `json:"player,omitempty"` // Optional, can be omitted if not needed
}

type Deck struct {
	ID           *uint    `json:"id,omitempty"`
	Commander    string   `json:"commander"`
	Crop         string   `json:"crop"`
	SecondaryImg string   `json:"secondary_image"`
	Image        string   `json:"image"`
	Colors       []string `json:"colors,omitempty"` // Scryfall color codes: W, U, B, R, G, C
	MoxfieldURL  *string  `json:"moxfield_url,omitempty"`
	Bracket      uint     `json:"bracket,omitempty"`
	Themes       []string `json:"themes,omitempty"`
}

type CreateDeckRequest struct {
	MoxfieldURL    *string  `json:"moxfield_url"`
	Themes         []string `json:"themes"`
	Bracket        uint     `json:"bracket"`
	Commander      string   `json:"commander"`
	Colors         []string `json:"colors"`
	Image          string   `json:"image"`
	SecondaryImage string   `json:"secondary_image"`
	Crop           string   `json:"crop"`
}

type Notification struct {
	ID               uint                 `json:"id"`
	Title            string               `json:"title"`
	Body             string               `json:"body"`
	Type             string               `json:"type"`
	Actions          []NotificationAction `json:"actions"`
	Read             bool                 `json:"read"`
	CreatedAt        time.Time            `json:"created_at"`
	GameID           *uint                `json:"game_id,omitempty"`
	ReferredPlayerID *string              `json:"referred_player_id,omitempty"`
	PlayerRankingID  *uint                `json:"player_ranking_id,omitempty"`
	Game             *Game                `json:"game,omitempty"`
	ReferredPlayer   *Player              `json:"referred_player,omitempty"`
}

type NotificationAction string

const (
	ActionDeleteRanking     NotificationAction = "delete_ranking"
	ActionViewGame          NotificationAction = "view_game"
	ActionAddImageGameEvent NotificationAction = "add_image_game_event"
)
