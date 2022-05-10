package deck

import "testing"

func TestNewFullDeck(t *testing.T) {
	deck := NewFullDeck()
	if len(deck) != 52 {
		t.Error("Expected the full deck to have 52 cards.")
	}

	suits := map[string][]string{}
	for _, card := range deck {
		suit := card[len(card)-1:]
		suits[suit] = append(suits[suit], card)
	}

	if len(suits) != 4 {
		t.Error("Expected the deck to have 4 suits")
	}
	for suit, cards := range suits {
		if len(cards) != 13 {
			t.Errorf("Expected the suit %s to have full range of cards (13) but got %d", suit, len(cards))
		}
	}
}

func TestShullfeDeck(t *testing.T) {
	deck := []*Card{}
	for _, card := range NewFullDeck()[:20] {
		deck = append(deck, &Card{
			Value: card,
		})
	}

	shuffled := append([]*Card{}, deck...)
	ShuffleDeck(shuffled)
	if compare(deck, shuffled) {
		t.Error("Expected the deck to be shuffled")
	}

	deck = []*Card{}
	for _, card := range NewFullDeck() {
		deck = append(deck, &Card{
			Value: card,
		})
	}
	shuffled = append([]*Card{}, deck...)
	ShuffleDeck(shuffled)
	if compare(deck, shuffled) {
		t.Error("Expected the full deck to be shuffled")
	}
}

func compare(deck1 []*Card, deck2 []*Card) bool {
	if len(deck1) != len(deck2) {
		return false
	}
	for i := 0; i < len(deck1); i++ {
		if deck1[i].Value != deck2[i].Value {
			return false
		}
	}
	return true
}

func TestValidateDeckCards(t *testing.T) {
	validDeck := []*Card{
		{
			Value: "AC",
		},
		{
			Value: "10H",
		},
		{
			Value: "JS",
		},
	}
	if err := ValidateDeckCards(validDeck); err != nil {
		t.Errorf("Expected the deck to be valid, but got a validation error instead: %s.", err.Error())
	}

	invalidDeck := []*Card{
		{
			Value: "C",
		},
		{
			Value: "10H",
		},
		{
			Value: "S",
		},
	}

	err := ValidateDeckCards(invalidDeck)
	if err == nil {
		t.Error("Expected the deck to be invalid.")
	}
	if err.Error() != "invalid cards values: C, S" {
		t.Errorf("Expected correct validation message, but got '%s' instead.", err.Error())
	}

	fullDeck := []*Card{}
	for _, card := range NewFullDeck() {
		fullDeck = append(fullDeck, &Card{
			Value: card,
		})
	}
	if err := ValidateDeckCards(fullDeck); err != nil {
		t.Errorf("Expected a full generated deck to be valid, but got error: '%s' instead.", err.Error())
	}
}

func TestAsCards(t *testing.T) {
	cards := AsCards("AC,   10C,5H   ")
	if len(cards) != 3 {
		t.Error("Expected to get 3 cards")
	}
	if !compare(cards, []*Card{
		{
			Value: "AC",
		},
		{
			Value: "10C",
		},
		{
			Value: "5H",
		},
	}) {
		t.Error("Expected to get the correct cards in order.")
	}
}
