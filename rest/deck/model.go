package deck

type CardResponse struct {
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}

type CreateDeckResponse struct {
	DeckID    string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}

type OpenDeckResponse struct {
	DeckID    string         `json:"deck_id"`
	Shuffled  bool           `json:"shuffled"`
	Remaining int            `json:"remaining"`
	Cards     []CardResponse `json:"cards"`
}

type DrawCardsResponse struct {
	Cards []CardResponse `json:"cards"`
}
