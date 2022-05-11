package deck

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/natemago/card-games-api/errors"
)

// Ranks is a list of all card ranks (A - K), not including Joker card.
var Ranks = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

// Suits is a list of card suits.
var Suits = []string{"C", "D", "H", "S"}

// RankNames is a mapping between the rank code and the rank name.
var RanksNames = map[string]string{
	"A":  "ACE",
	"2":  "2",
	"3":  "3",
	"4":  "4",
	"5":  "5",
	"6":  "6",
	"7":  "7",
	"8":  "8",
	"9":  "9",
	"10": "10",
	"J":  "JACK",
	"Q":  "QUEEN",
	"K":  "KING",
}

// SuitsNames is a mapping between the suit code (short 1 letter) to the actual name.
var SuitsNames = map[string]string{
	"C": "CLUBS",
	"D": "DIAMONDS",
	"H": "HEARTS",
	"S": "SPADES",
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewFullDeck returns a new full deck of 52 cards, sorted in order.
func NewFullDeck() []string {
	deck := []string{}

	for _, suit := range Suits {
		for _, rank := range Ranks {
			deck = append(deck, fmt.Sprintf("%s%s", rank, suit))
		}
	}

	return deck
}

// ShuffleDeck shuffles a deck of cards. The deck does not have to be full.
func ShuffleDeck(deck []*Card) {
	rand.Shuffle(len(deck), func(i, j int) {
		t := deck[i]
		deck[i] = deck[j]
		deck[j] = t
	})
	for i, card := range deck {
		card.Idx = i
	}
}

// ValidateDeckCards validates if the cards are actually valid cards and there are no duplicates in the deck.
func ValidateDeckCards(cards []*Card) error {
	var invalidCards []string
	seen := map[string]bool{}
	for _, card := range cards {
		if card.RankName() == "" || card.SuitName() == "" {
			invalidCards = append(invalidCards, card.Value)
			continue
		}
		if _, ok := seen[card.Value]; ok {
			// duplicate
			invalidCards = append(invalidCards, card.Value)
		}
		seen[card.Value] = true
	}

	if len(invalidCards) > 0 {
		return errors.ValidationError(fmt.Sprintf("invalid cards values: %s", strings.Join(invalidCards, ", ")), nil)
	}

	return nil
}

// AsCards parses a string of comma separated values into a deck of cards.
// Note that after parsing, some of the cards may hold invalid value or the deck might
// have duplicate cards. See ValidateDeckCards to validate the deck.
func AsCards(cardsStr string) []*Card {
	var result []*Card

	for _, cardValue := range strings.Split(strings.TrimSpace(cardsStr), ",") {
		cardValue = strings.TrimSpace(cardValue)
		if cardValue == "" {
			continue
		}
		result = append(result, &Card{
			Value: cardValue,
		})
	}

	return result
}
