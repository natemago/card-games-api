package deck

import (
	"testing"

	"github.com/natemago/card-games-api/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestData struct {
	DB            *gorm.DB
	FullDeckID    string
	PartialDeckID string
}

func setupTest(t *testing.T) (TestData, func(*testing.T)) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Errorf("Failed to setup database: %s", err.Error())
	}

	AutoMigrateDeckModels(db)

	deckRepo := NewDBDeckRepository(db)

	partialDeck, err := deckRepo.CreateDeck(&Deck{
		Cards: []*Card{
			{
				Value: "2C",
			},
			{
				Value: "5H",
			},
			{
				Value: "KD",
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create partial deck: %s", err.Error())
	}

	fullDeck, err := deckRepo.CreateDeck(&Deck{})
	if err != nil {
		t.Fatalf("Failed to create full deck: %s", err.Error())
	}

	return TestData{
			DB:            db,
			PartialDeckID: partialDeck.ID,
			FullDeckID:    fullDeck.ID,
		}, func(tb *testing.T) {

		}
}

func TestCreateDeck_FullDeck(t *testing.T) {
	td, tearDown := setupTest(t)
	defer tearDown(t)

	deckRepo := NewDBDeckRepository(td.DB)

	result, err := deckRepo.CreateDeck(&Deck{})
	if err != nil {
		t.Fatalf("Expected to create a new full deck, but got error: %s", err.Error())
	}
	if result.Shuffled {
		t.Errorf("Expected the deck to not be shuffled.")
	}
	if result.Remaining != 52 {
		t.Error("Expected the deck to have 52 cards remaining.")
	}
	if len(result.Cards) != 52 {
		t.Error("Expected the deck to have 52 actual cards.")
	}
	if !cardsInOrder(result.Cards) {
		t.Error("Expected all cards to be in proper order.")
	}

	result, err = deckRepo.CreateDeck(&Deck{
		Shuffled: true,
	})
	if err != nil {
		t.Fatalf("Expected to create a new full deck, but got error: %s", err.Error())
	}
	if !result.Shuffled {
		t.Errorf("Expected the deck to be shuffled.")
	}
	if result.Remaining != 52 {
		t.Error("Expected the deck to have 52 cards remaining.")
	}
	if len(result.Cards) != 52 {
		t.Error("Expected the deck to have 52 actual cards.")
	}
	if cardsInOrder(result.Cards) {
		t.Error("Expected the cards to be shuffled.")
	}
}

func TestCreateDeck_PartialDeck(t *testing.T) {
	td, tearDown := setupTest(t)
	defer tearDown(t)

	deckRepo := NewDBDeckRepository(td.DB)

	result, err := deckRepo.CreateDeck(&Deck{
		Cards: []*Card{
			{
				Value: "2C",
			},
			{
				Value: "4H",
			},
			{
				Value: "KD",
			},
		},
	})
	if err != nil {
		t.Fatalf("Expected to create a partial deck, but got error instead: %s", err.Error())
	}
	if result.Shuffled {
		t.Errorf("Expected the deck to not be shuffled.")
	}
	if result.Remaining != 3 {
		t.Error("Expected the deck to have 3 cards remaining.")
	}
	if len(result.Cards) != 3 {
		t.Error("Expected the deck to have 3 actual cards.")
	}

	suppliedCards := []*Card{}
	for _, card := range NewFullDeck()[:20] {
		suppliedCards = append(suppliedCards, &Card{
			Value: card,
		})
	}

	result, err = deckRepo.CreateDeck(&Deck{
		Shuffled: true,
		Cards:    append([]*Card{}, suppliedCards...),
	})
	if err != nil {
		t.Fatalf("Expected to create a partial deck, but got error instead: %s", err.Error())
	}
	if !result.Shuffled {
		t.Errorf("Expected the deck to be shuffled.")
	}
	if result.Remaining != 20 {
		t.Error("Expected the deck to have 20 cards remaining.")
	}
	if len(result.Cards) != 20 {
		t.Error("Expected the deck to have 20 actual cards.")
	}
	if compare(suppliedCards, result.Cards) {
		t.Error("Expected the deck to be shuffled.")
	}
}

func TestCreateDeck_InvaidCards(t *testing.T) {
	td, tearDown := setupTest(t)
	defer tearDown(t)

	deckRepo := NewDBDeckRepository(td.DB)

	_, err := deckRepo.CreateDeck(&Deck{
		Cards: []*Card{
			{
				Value: "2C",
			},
			{
				Value: "TA",
			},
			{
				Value: "KD",
			},
		},
	})
	if err == nil {
		t.Fatalf("Expected to get validation error.")
	}
	if !errors.IsValidationError(err) {
		t.Error("Expected the actual error to be ValidationError.")
	}
}

func TestGetDeck(t *testing.T) {
	td, tearDown := setupTest(t)
	defer tearDown(t)

	deckRepo := NewDBDeckRepository(td.DB)

	result, err := deckRepo.GetDeck(td.FullDeckID)
	if err != nil {
		t.Fatalf("Expected to get the full deck of cards, but got error instead: %s", err.Error())
	}
	if len(result.Cards) != 52 {
		t.Error("Expected to get the full deck.")
	}

	result, err = deckRepo.GetDeck(td.PartialDeckID)
	if err != nil {
		t.Fatalf("Expected to get the partial deck of cards, but got error instead: %s", err.Error())
	}
	if len(result.Cards) != 3 {
		t.Error("Expected to get the partial deck.")
	}
}

func TestGetDeck_NotFound(t *testing.T) {
	td, tearDown := setupTest(t)
	defer tearDown(t)

	deckRepo := NewDBDeckRepository(td.DB)

	_, err := deckRepo.GetDeck("00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Error("Expected not to find the non-existend deck.")
	}
	if !errors.IsNotFoundError(err) {
		t.Error("Expected the error to be NotFoundError.")
	}
}

func TestDrawCards(t *testing.T) {
	td, tearDown := setupTest(t)
	defer tearDown(t)

	deckRepo := NewDBDeckRepository(td.DB)

	result, err := deckRepo.DrawCards(td.FullDeckID, 1)
	if err != nil {
		t.Fatalf("Expected to draw 1 card, but got an error instead: %s", err.Error())
	}
	if len(result) != 1 {
		t.Fatalf("Expected to draw exactly 1 card, but actually drawn: %d", len(result))
	}
	if result[0].Value != "AC" {
		t.Error("Expected the drawn card to be 'AC'.")
	}

	deck, err := deckRepo.GetDeck(td.FullDeckID)
	if err != nil {
		t.Fatal("Expected to get the deck back.")
	}
	if deck.Remaining != 51 {
		t.Errorf("Expected the deck to have 51 remaining cards, but it has: %d", deck.Remaining)
	}

	result, err = deckRepo.DrawCards(td.FullDeckID, 10)
	if err != nil {
		t.Fatalf("Expected to draw 1 card, but got an error instead: %s", err.Error())
	}
	if len(result) != 10 {
		t.Errorf("Expected to draw exactly 10 cards, but actually drawn: %d", len(result))
	}

	deck, err = deckRepo.GetDeck(td.FullDeckID)
	if err != nil {
		t.Fatal("Expected to get the deck back.")
	}
	if deck.Remaining != 41 {
		t.Errorf("Expected the deck to have 41 remaining cards, but it has: %d", deck.Remaining)
	}
}

func TestDrawCards_Overdraw(t *testing.T) {
	td, tearDown := setupTest(t)
	defer tearDown(t)

	deckRepo := NewDBDeckRepository(td.DB)

	_, err := deckRepo.DrawCards(td.PartialDeckID, 10)
	if err == nil {
		t.Fatal("Expected to get an overdraw error.")
	}
	if !errors.IsBadRequestError(err) {
		t.Error("Expected the error to be BadRequestError.")
	}
}

func cardsInOrder(cards []*Card) bool {
	deckInOrder := NewFullDeck()
	for i, card := range cards {
		if deckInOrder[i] != card.Value {
			return false
		}
	}
	return true
}
