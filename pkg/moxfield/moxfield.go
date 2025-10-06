package moxfield

const (
	moxFieldSearchUrl = "https://api2.moxfield.com/v2/decks/search-sfw"
	moxFieldDeckUrl   = "https://api2.moxfield.com/v3/decks/all/"
)

type Deck struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	MoxfieldID string   `json:"moxfield_id"`
	Themes     []string `json:"themes"`
	Bracket    uint     `json:"bracket"`
	Commander  string   `json:"commander"`
	Colors     []string `json:"colors"`
	//Image          string   `json:"image"`
	//SecondaryImage *string  `json:"secondary_image"`
	//Crop           string   `json:"crop"`
}

func GetDecksForUser(username string) ([]Deck, error) {
	return []Deck{}, nil
}

func GetDeckByID(deckID string) (*Deck, error) {
	return &Deck{}, nil
}
