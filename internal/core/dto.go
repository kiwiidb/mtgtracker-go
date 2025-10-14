package core

import (
	"time"
)

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
	RankingID      uint             `json:"ranking_id"`
	Position       int              `json:"position"`
	Description    *GameDescription `json:"description,omitempty"`
	StartingPlayer *bool            `json:"starting_player,omitempty"`
}

type UpdateRankingRequest struct {
	Description    *GameDescription `json:"description,omitempty"`
	StartingPlayer *bool            `json:"starting_player,omitempty"`
	CouldHaveWon   *bool            `json:"could_have_won,omitempty"`
	EarlySolRing   *bool            `json:"early_sol_ring,omitempty"`
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
	Player PlayerResponse `json:"player"`
	Count  int            `json:"count"`
}
type DeckWithCount struct {
	Deck  DeckResponse `json:"deck"`
	Count int          `json:"count"`
	Wins  int          `json:"wins"`
}
type PlayerResponse struct {
	ID                   string          `json:"id"`
	Name                 string          `json:"name"`
	ProfileImageURL      string          `json:"profile_image_url,omitempty"`
	MoxfieldUsername     string          `json:"moxfield_username,omitempty"`
	Colors               []string        `json:"colors,omitempty"` // Top 2 most played colors
	WinrateAllTime       float64         `json:"winrate_all_time"`
	NumberofGamesAllTime int             `json:"number_of_games_all_time"`
	DecksAllTime         []DeckWithCount `json:"decks_all_time"`
	CurrentGame          *GameResponse   `json:"current_game,omitempty"`
}

type UpdatePlayerRequest struct {
	MoxfieldUsername *string `json:"moxfield_username,omitempty"`
}

type GameResponse struct {
	ID         uint                `json:"id"`
	CreatorID  *string             `json:"creator_id,omitempty"`
	Duration   *int                `json:"duration,omitempty"`
	Date       *time.Time          `json:"date,omitempty"`
	EndDate    *time.Time          `json:"end_date,omitempty"`
	Comments   string              `json:"comments,omitempty"`
	Rankings   []RankingResponse   `json:"rankings,omitempty"`
	Finished   bool                `json:"finished"`
	GameEvents []GameEventResponse `json:"game_events,omitempty"`
	Creator    *PlayerResponse     `json:"creator,omitempty"`
}

type GameEventResponse struct {
	GameID               uint             `json:"game_id"`
	EventType            string           `json:"event_type"`
	DamageDelta          int              `json:"damage_delta"`
	CreatedAt            time.Time        `json:"created_at"`
	TargetLifeTotalAfter int              `json:"target_life_total_after"`
	SourceRanking        *RankingResponse `json:"source_ranking,omitempty"`
	TargetRanking        *RankingResponse `json:"target_ranking,omitempty"`
	ImageUrl             string           `json:"image_url"`                  // URL of the uploaded image
	UploadImageUrl       string           `json:"upload_image_url,omitempty"` // Presigned URL for image upload
	Comment              *string          `json:"comment,omitempty"`          // New field for text description
}

type RankingResponse struct {
	ID                     uint             `json:"id"`
	PlayerID               *string          `json:"player_id,omitempty"`
	Position               int              `json:"position"`
	LastLifeTotal          *int             `json:"last_life_total,omitempty"`
	LastLifeTotalTimestamp *time.Time       `json:"last_life_total_timestamp,omitempty"`
	Deck                   DeckResponse     `json:"deck"`
	Player                 *PlayerResponse  `json:"player,omitempty"` // Optional, can be omitted if not needed
	Description            *GameDescription `json:"description,omitempty"`
}

type DeckResponse struct {
	ID           *uint    `json:"id,omitempty"`
	Commander    string   `json:"commander"`
	Crop         string   `json:"crop"`
	SecondaryImg string   `json:"secondary_image"`
	Image        string   `json:"image"`
	Colors       []string `json:"colors,omitempty"` // Scryfall color codes: W, U, B, R, G, C
	MoxfieldURL  *string  `json:"moxfield_url,omitempty"`
	Bracket      *uint    `json:"bracket,omitempty"`
	Themes       []string `json:"themes,omitempty"`
}

type CreateDeckRequest struct {
	MoxfieldURL    *string  `json:"moxfield_url"`
	Themes         []string `json:"themes"`
	Bracket        *uint    `json:"bracket,omitempty"`
	Commander      string   `json:"commander"`
	Colors         []string `json:"colors"`
	Image          string   `json:"image"`
	SecondaryImage string   `json:"secondary_image"`
	Crop           string   `json:"crop"`
}

type SearchGamesRequest struct {
	PlayerIDs     []string `json:"player_ids,omitempty"`      // Games where ANY of these players participated (OR)
	Commanders    []string `json:"commanders,omitempty"`      // Games where ANY of these commanders were played (OR)
	AllPlayers    []string `json:"all_players,omitempty"`     // Games where ALL of these players participated (AND)
	AllCommanders []string `json:"all_commanders,omitempty"`  // Games where ALL of these commanders were played (AND)
}

func (req SearchGamesRequest) ToFilter() GameFilter {
	return GameFilter{
		PlayerIDs:     req.PlayerIDs,
		Commanders:    req.Commanders,
		AllPlayers:    req.AllPlayers,
		AllCommanders: req.AllCommanders,
	}
}
