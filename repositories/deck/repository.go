package deck

type DeckRepository interface {
	CreateDeck(deck *Deck) (*Deck, error)
	GetDeck(deckID string) (*Deck, error)
	DrawCards(deckID string, numCards int) ([]*Card, error)
}
