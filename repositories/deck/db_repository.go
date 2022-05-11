package deck

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	api_errors "github.com/natemago/card-games-api/errors"
)

// DBDeckRepository holds the reference to the underlying database connection
// and binds method for DeckRepository interface.
type DBDeckRepository struct {
	db *gorm.DB
}

// CreateDeck creates a new deck of cards.
// To generate a full 52 deck of cards in order, supply a pointer to an empty Deck struct.
// By setting Deck.Shuffled to true, it will generate a shuffled deck,
// Setting Deck.Cards to a non-empty array of Card, it will generate a deck with the supplied cards only.
// If cards are supplied, it may return a ValidationError if some of the cards have multiple values or are
// duplicates.
func (d *DBDeckRepository) CreateDeck(deck *Deck) (*Deck, error) {
	if deck.ID == "" {
		deck.ID = uuid.New().String()
	}

	if deck.Cards == nil {
		cards := NewFullDeck()
		for i, card := range cards {
			deck.Cards = append(deck.Cards, &Card{
				DeckID: deck.ID,
				Value:  card,
				Idx:    i,
			})
		}
	}

	if err := ValidateDeckCards(deck.Cards); err != nil {
		return nil, err
	}

	if deck.Shuffled {
		ShuffleDeck(deck.Cards)
	}

	deck.Remaining = len(deck.Cards)

	if err := d.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&deck)
		if result.Error != nil {
			return result.Error
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return deck, nil
}

// GetDeck looks up a Deck by its ID and returns a reference to a populated Deck.
// If there is no deck with the given ID, then a NotFound error is returned.
func (d *DBDeckRepository) GetDeck(deckID string) (*Deck, error) {
	deck := &Deck{}

	result := d.db.Where("id=?", deckID).First(deck)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, api_errors.NotFoundError("no such deck", nil)
		}
		return nil, result.Error
	}

	cards := []*Card{}

	result = d.db.Debug().Where("deck_id = ? AND drawn IS FALSE", deckID).Order("idx").Find(&cards)

	if result.Error != nil {
		return nil, result.Error
	}

	deck.Cards = cards

	return deck, nil
}

// DrawCards draws a number of cards from the deck.
// Once drawn, the cards will no longer be in the deck.
// Returns a list of the drawn cards.
// If there is no deck with the given deckID, then a NotFoundError will be returned.
// If the number of cards is less then one, then just one card will be returned.
// If the number of cards to be drawn is greater than the number of remaining cards in the deck,
// then a BadRequestError will be returned.
// After the cards are drawn, the remaning number of cards in the deck will decrease by the number of
// drawn cards.
func (d *DBDeckRepository) DrawCards(deckID string, numCards int) ([]*Card, error) {
	var drawn []*Card
	if err := d.db.Transaction(func(tx *gorm.DB) error {
		deck, err := d.GetDeck(deckID)
		if err != nil {
			return err
		}

		if numCards < 1 {
			numCards = 1
		}

		if numCards > len(deck.Cards) {
			return api_errors.BadRequestError("not enough cards in deck", nil)
		}

		drawn = deck.Cards[0:numCards]
		for _, card := range drawn {
			card.Drawn = true
			result := d.db.Save(card)
			if result.Error != nil {
				return result.Error
			}
		}

		deck.Remaining -= numCards

		result := d.db.Save(deck)
		if result.Error != nil {
			return result.Error
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return drawn, nil
}

// NewDBDeckRepository creates a new DeckRepository with the given database connection.
func NewDBDeckRepository(db *gorm.DB) DeckRepository {
	return &DBDeckRepository{
		db: db,
	}
}

// AutoMigrateDeckModels performs an automatic migration of the defined Gorm models in the database.
func AutoMigrateDeckModels(db *gorm.DB) error {
	if err := db.AutoMigrate(&Deck{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Card{}); err != nil {
		return err
	}

	return nil
}
