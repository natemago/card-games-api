package deck

// DeckRepository defines methods for managing a deck of cards, like creating, showing the deck or drawing a card from it.
type DeckRepository interface {

	// CreateDeck creates new deck of cards.
	// To generate a full 52 deck of cards in order, supply a pointer to an empty Deck struct.
	// By setting Deck.Shuffled to true, it will generate a shuffled deck,
	// Setting Deck.Cards to a non-empty array of Card, it will generate a deck with the supplied cards only.
	// If cards are supplied, it may return a ValidationError if some of the cards have multiple values or are
	// duplicates.
	CreateDeck(deck *Deck) (*Deck, error)

	// GetDeck looks up a deck of cards by its ID.
	// If there is no deck with the given ID, then a NotFound error is returned.
	GetDeck(deckID string) (*Deck, error)

	// DrawCards draws a number of cards from the deck.
	// Once drawn, the cards will no longer be in the deck.
	// Returns a list of the drawn cards.
	// If there is no deck with the given deckID, then a NotFoundError will be returned.
	// If the number of cards is less then one, then just one card will be returned.
	// If the number of cards to be drawn is greater than the number of remaining cards in the deck,
	// then a BadRequestError will be returned.
	// After the cards are drawn, the remaining number of cards in the deck will decrease by the number of
	// drawn cards.
	DrawCards(deckID string, numCards int) ([]*Card, error)
}
