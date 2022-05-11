package deck

// CardResponse represents a Card response object. Holds the data for a particular card in a deck.
type CardResponse struct {
	// Value is the card rank (number), like "ACE", "2", "10", "QUEEN" etc.
	Value string `json:"value"`

	// Suit is the card suit name, like "HEARTS" or "DIAMONDS".
	Suit string `json:"suit"`

	// Code is the full card code, like: "AC" (Ace of Clubs), "2H" (Two of hearts) etc.
	Code string `json:"code"`
}

// CreateDeckResponse represents the response of a CreateDeck call.
type CreateDeckResponse struct {
	// DeckID is the generated deck id for the new deck.
	DeckID string `json:"deck_id"`

	// Shuffled flag whether the deck is shuffled or in proper order.
	Shuffled bool `json:"shuffled"`

	// Remaining is the number of remaining cards in the deck.
	Remaining int `json:"remaining"`
}

// OpenDeckResponse represents the response for an OpenDeck call (show all cards in deck).
type OpenDeckResponse struct {
	// DeckID is the id of the deck.
	DeckID string `json:"deck_id"`
	// Shuffled flag whether the deck is shuffled or in proper order.
	Shuffled bool `json:"shuffled"`
	// Remaining is the number of remaining cards in the deck.
	Remaining int `json:"remaining"`

	// Cards is the list of cards in the deck, in the order they were inserted/generated.
	Cards []CardResponse `json:"cards"`
}

// DrawCardsResponse represents the response for a DrawCards call - draw one or more cards.
// Holds the list of the drawn cards.
type DrawCardsResponse struct {
	// Cards list of cards drawn from the deck.
	Cards []CardResponse `json:"cards"`
}
