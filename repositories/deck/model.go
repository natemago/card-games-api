package deck

import (
	"time"
)

// Deck represents the database model for a deck of cards.
type Deck struct {
	// ID is a unique identifier for this deck, usually an UUID v4.
	ID string `gorm:"primaryKey"`

	// CreatedAt is the time when this deck was created.
	CreatedAt time.Time

	// UpdatedAt is the time when this deck was last updated.
	UpdatedAt time.Time

	// Shuffled flag - whether this deck is shuffled.
	Shuffled bool

	// Remaining is the number of remaining cards in the deck.
	Remaining int

	// Cards is the list of actual cards, in the given order (proper or shuffled) in the deck.
	Cards []*Card
}

// Card represents the database model for a particular card belonging to a deck.
type Card struct {
	// DeckID is the foreign key to the parent deck of cards.
	DeckID string `gorm:"primaryKey"`

	// Value is the actual value of the card (code). For example: "AC", "10S", "KH" etc.
	Value string `gorm:"primaryKey"`

	// Drawn is a flag whether this card was drawn or not.
	Drawn bool

	// Idx is the index of the card used to determine the order of the cards in the particular deck.
	Idx int
}

// SuitName returns the suit of the card. For example for the card "KH" it will return "HEARTS".
func (c *Card) SuitName() string {
	if c.Value == "" {
		return ""
	}
	suit := c.Value[len(c.Value)-1:]
	suitName, ok := SuitsNames[suit]
	if !ok {
		return ""
	}
	return suitName
}

// RankName returns the card rank (number). For example the card "QS", the rank is "QUEEN".
func (c *Card) RankName() string {
	if c.Value == "" {
		return ""
	}
	rank := c.Value[:len(c.Value)-1]
	rankName, ok := RanksNames[rank]
	if !ok {
		return ""
	}
	return rankName
}
