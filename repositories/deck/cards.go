package deck

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"toggl.com/services/card-games-api/errors"
)

var Ranks = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
var Suits = []string{"C", "D", "H", "S"}

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
var SuitsNames = map[string]string{
	"C": "CLUBS",
	"D": "DIAMONDS",
	"H": "HEARTS",
	"S": "SPADES",
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewFullDeck() []string {
	deck := []string{}

	for _, suit := range Suits {
		for _, rank := range Ranks {
			deck = append(deck, fmt.Sprintf("%s%s", rank, suit))
		}
	}

	return deck
}

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

func ValidateDeckCards(cards []*Card) error {
	var invalidCards []string
	for _, card := range cards {
		if card.RankName() == "" || card.SuitName() == "" {
			invalidCards = append(invalidCards, card.Value)
		}
	}

	if len(invalidCards) > 0 {
		return errors.ValidationError(fmt.Sprintf("invalid cards values: %s", strings.Join(invalidCards, ", ")), nil)
	}

	return nil
}

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
