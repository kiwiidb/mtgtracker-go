package moxfield

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	moxFieldSearchUrl = "https://api2.moxfield.com/v2/decks/search-sfw"
	moxFieldDeckUrl   = "https://api2.moxfield.com/v3/decks/all/"
)

type Deck struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	MoxfieldID     string   `json:"moxfield_id"`
	Themes         []string `json:"themes"`
	Bracket        uint     `json:"bracket"`
	Commander      string   `json:"commander"`
	Colors         []string `json:"colors"`
	Image          string   `json:"image"`
	SecondaryImage *string  `json:"secondary_image"`
	Crop           string   `json:"crop"`
}

type searchResponse struct {
	Data []searchDeck `json:"data"`
}

type searchDeck struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	PublicID      string   `json:"publicId"`
	Format        string   `json:"format"`
	ColorIdentity []string `json:"colorIdentity"`
	Bracket       uint     `json:"bracket"`
}

type deckResponse struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	ColorIdentity []string `json:"colorIdentity"`
	Bracket       uint     `json:"bracket"`
	Hubs          []struct {
		Name string `json:"name"`
	} `json:"hubs"`
	Main struct {
		ScryfallID string `json:"scryfall_id"`
		CardFaces  []struct {
			Name string `json:"name"`
		} `json:"card_faces"`
	} `json:"main"`
	Boards struct {
		Commanders struct {
			Cards map[string]struct {
				Card struct {
					Name       string `json:"name"`
					ScryfallID string `json:"scryfall_id"`
					CardFaces  []struct {
						Name string `json:"name"`
					} `json:"card_faces"`
				} `json:"card"`
			} `json:"cards"`
		} `json:"commanders"`
	} `json:"boards"`
}

// buildScryfallImageURL constructs a Scryfall image URL from a scryfall_id
// face can be "front" or "back"
func buildScryfallImageURL(scryfallID, imageType, face string) string {
	if len(scryfallID) < 2 {
		return ""
	}
	// Extract first two characters for path segments
	dir1 := string(scryfallID[0])
	dir2 := string(scryfallID[1])

	return fmt.Sprintf("https://cards.scryfall.io/%s/%s/%s/%s/%s.jpg", imageType, face, dir1, dir2, scryfallID)
}

func GetDecksForUser(username string) ([]Deck, error) {
	// Build search URL with username parameter
	searchURL := fmt.Sprintf("%s?authorUserName=%s&pageSize=100&fmt=commander", moxFieldSearchUrl, url.QueryEscape(username))

	// Make HTTP request with User-Agent header
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "MTGTracker/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch decks: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("moxfield API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse search response
	var searchResp searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	// Fetch detailed info for each deck to get commander name
	decks := make([]Deck, 0, len(searchResp.Data))
	for _, deckInfo := range searchResp.Data {
		// Only include commander format decks
		if deckInfo.Format != "commander" {
			continue
		}

		deck, err := GetDeckByID(deckInfo.PublicID)
		if err != nil {
			// Log error but continue with other decks
			fmt.Printf("Warning: failed to fetch deck %s: %v\n", deckInfo.ID, err)
			continue
		}

		// Use bracket from search response
		deck.Bracket = deckInfo.Bracket
		decks = append(decks, *deck)
	}

	return decks, nil
}

func GetDeckByID(deckID string) (*Deck, error) {
	// Fetch deck details
	deckURL := fmt.Sprintf("%s%s", moxFieldDeckUrl, url.PathEscape(deckID))

	// Make HTTP request with User-Agent header
	req, err := http.NewRequest("GET", deckURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "MTGTracker/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deck: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("moxfield API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse deck response
	var deckResp deckResponse
	if err := json.NewDecoder(resp.Body).Decode(&deckResp); err != nil {
		return nil, fmt.Errorf("failed to decode deck response: %w", err)
	}

	// Extract themes from hubs
	themes := make([]string, 0, len(deckResp.Hubs))
	for _, hub := range deckResp.Hubs {
		themes = append(themes, hub.Name)
	}

	// Extract commander name and scryfall_id from commanders board
	var commanderName, scryfallID string
	var hasCardFaces bool
	for _, commanderCard := range deckResp.Boards.Commanders.Cards {
		commanderName = commanderCard.Card.Name
		scryfallID = commanderCard.Card.ScryfallID
		hasCardFaces = len(commanderCard.Card.CardFaces) > 0
		break // Take the first commander
	}

	// Build image URLs from scryfall_id
	imageURL := buildScryfallImageURL(scryfallID, "normal", "front")
	cropURL := buildScryfallImageURL(scryfallID, "art_crop", "front")

	// If card has multiple faces, secondary image is the back face
	var secondaryImageURL *string
	if hasCardFaces {
		backImageURL := buildScryfallImageURL(scryfallID, "normal", "back")
		secondaryImageURL = &backImageURL
	}

	return &Deck{
		ID:             deckResp.ID,
		Name:           deckResp.Name,
		MoxfieldID:     deckID,
		Commander:      commanderName,
		Colors:         deckResp.ColorIdentity,
		Themes:         themes,
		Bracket:        deckResp.Bracket,
		Image:          imageURL,
		SecondaryImage: secondaryImageURL,
		Crop:           cropURL,
	}, nil
}
