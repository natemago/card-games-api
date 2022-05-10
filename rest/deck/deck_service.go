package deck

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	deck_repo "github.com/natemago/card-games-api/repositories/deck"
)

type DeckService struct {
	Repository deck_repo.DeckRepository
}

func (d *DeckService) CreateDeck(ctx *gin.Context) {
	cardsParam, _ := ctx.GetQuery("cards")
	_, shuffled := ctx.GetQuery("shuffled")

	var cards []*deck_repo.Card

	if cardsParam != "" {
		cards = deck_repo.AsCards(cardsParam)
	}

	deck, err := d.Repository.CreateDeck(&deck_repo.Deck{
		Shuffled: shuffled,
		Cards:    cards,
	})

	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, &CreateDeckResponse{
		DeckID:    deck.ID,
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining,
	})
}

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

func NewDeckService(deckRepository deck_repo.DeckRepository) *DeckService {
	return &DeckService{
		Repository: deckRepository,
	}
}
