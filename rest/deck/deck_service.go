package deck

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	deck_repo "github.com/natemago/card-games-api/repositories/deck"
)

// DeckService represents the REST API service for the Deck resource.
// Uses the DeckRepository to actually manage the deck of cards.
type DeckService struct {
	Repository deck_repo.DeckRepository
}

// CreateDeck endpoint for creating new deck given.
// Accepts two query parameters:
//  - shuffled - (optional) whether to create a shuffled deck or a deck with the cards in proper order.
//  - cards - (optional) an optional list of cards given in a comma-separated string. When supplied, the
//      deck will contain only the given cards (partial deck).
// If none of the query parameters are supplied, then a full 52 deck of cards in proper order will be created.
// If the cards list contain any invalid or duplicated values, returns a 400 Bad Request error response.
func (d *DeckService) CreateDeck(ctx *gin.Context) {
	cardsParam, _ := ctx.GetQuery("cards")
	shuffledParam, _ := ctx.GetQuery("shuffled")

	var cards []*deck_repo.Card
	shuffled := false

	if cardsParam != "" {
		cards = deck_repo.AsCards(cardsParam)
	}

	if shuffledParam != "" {
		shuffled, _ = strconv.ParseBool(strings.TrimSpace(shuffledParam))
	}

	deck, err := d.Repository.CreateDeck(&deck_repo.Deck{
		Shuffled: shuffled,
		Cards:    cards,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, &CreateDeckResponse{
		DeckID:    deck.ID,
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining,
	})
}

// OpenDeck looks up a deck by its id, and returns the deck data.
// Accepts one path parameter: deckId - the ID of the deck to look up.
// If the deck does not exist, generates a 404 error response.
// Returns the deck metadata and the list of cards remaining in the deck.
func (d *DeckService) OpenDeck(ctx *gin.Context) {
	deckID := ctx.Param("deckId")
	if deckID == "" {
		ctx.Error(fmt.Errorf("not-found"))
		return
	}

	deck, err := d.Repository.GetDeck(deckID)
	if err != nil {
		ctx.Error(err)
		return
	}

	var cards []CardResponse

	for _, card := range deck.Cards {
		cards = append(cards, CardResponse{
			Value: card.RankName(),
			Suit:  card.SuitName(),
			Code:  card.Value,
		})
	}

	ctx.JSON(http.StatusOK, &OpenDeckResponse{
		DeckID:    deck.ID,
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining,
		Cards:     cards,
	})
}

// DrawCards draws a number of cards from a given deck.
// Accepts two parameters:
//  - deckId - a path parameter. The ID of the deck to draw cards from.
//  - count - query parameter, integer. The number of cards to draw from the deck.
// Returns a list of the drawn cards.
// If there is no deck with the given id, then returns a 404 not found error response.
// If the count paramters is not an integer or is greater then the number of remaining cards,
// then returns a 400 Bad Request error response.
func (d *DeckService) DrawCards(ctx *gin.Context) {
	deckID := ctx.Param("deckId")
	if deckID == "" {
		ctx.Error(fmt.Errorf("not-found"))
		return
	}

	numCards := 1

	if numCardsStr, ok := ctx.GetQuery("count"); ok {
		numCardsStr = strings.TrimSpace(numCardsStr)
		if numCardsStr != "" {
			var err error
			numCards, err = strconv.Atoi(numCardsStr)
			if err != nil {
				ctx.Error(fmt.Errorf("invalid-cards-count-number"))
				return
			}
		}
	}

	if numCards < 1 {
		ctx.Error(fmt.Errorf("invalid-cards-count-number"))
		return
	}

	drawnCards, err := d.Repository.DrawCards(deckID, numCards)
	if err != nil {
		ctx.Error(err)
		return
	}

	var respCards []CardResponse

	for _, card := range drawnCards {
		respCards = append(respCards, CardResponse{
			Code:  card.Value,
			Suit:  card.SuitName(),
			Value: card.RankName(),
		})
	}

	ctx.JSON(http.StatusOK, &DrawCardsResponse{
		Cards: respCards,
	})
}

// NewDeckService creates a new pointer to a DeckService using the given DeckRepository.
func NewDeckService(deckRepository deck_repo.DeckRepository) *DeckService {
	return &DeckService{
		Repository: deckRepository,
	}
}
