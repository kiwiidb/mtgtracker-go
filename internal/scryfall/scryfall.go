package scryfall

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// CardImageURIs holds the image URIs returned by Scryfall for a card
type CardImageURIs struct {
	Small      string `json:"small"`
	Normal     string `json:"normal"`
	Large      string `json:"large"`
	Png        string `json:"png"`
	ArtCrop    string `json:"art_crop"`
	BorderCrop string `json:"border_crop"`
}

// Card represents the part of Scryfall's response that includes image URIs
type Card struct {
	Name       string        `json:"name"`
	ImageURIs  CardImageURIs `json:"image_uris"`
	CardFaces  []Card        `json:"card_faces,omitempty"` // Optional for cards with multiple faces
	OracleText string        `json:"oracle_text"`
	Power      string        `json:"power"`
	Toughness  string        `json:"toughness"`
	Colors     []string      `json:"colors"`
	ManaCost   string        `json:"mana_cost"`
	TypeLine   string        `json:"type_line"`
}

func GetCard(cardName string) (*Card, error) {
	baseURL := "https://api.scryfall.com/cards/named"
	params := url.Values{}
	params.Add("exact", cardName)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch card data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Scryfall API returned status: %s", resp.Status)
	}

	var card Card
	if err := json.NewDecoder(resp.Body).Decode(&card); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &card, nil
}
