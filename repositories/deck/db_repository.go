package deck

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	api_errors "github.com/natemago/card-games-api/errors"
)

type DBDeckRepository struct {
	db *gorm.DB
}

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

func NewDBDeckRepository(db *gorm.DB) DeckRepository {
	return &DBDeckRepository{
		db: db,
	}
}

func AutoMigrateDeckModels(db *gorm.DB) error {
	if err := db.AutoMigrate(&Deck{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Card{}); err != nil {
		return err
	}

	return nil
}
