package deck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/natemago/card-games-api/config"
	"github.com/natemago/card-games-api/errors"
	"github.com/natemago/card-games-api/repositories"
	deck_repo "github.com/natemago/card-games-api/repositories/deck"
)

var testConfig = &config.Config{
	APIConfig: config.APIConfig{
		Host: "",
		Port: 8080,
	},
	DBConfig: config.DBConfig{
		Dialect: "sqlite",
		URL:     "file::memory:?cache=shared",
	},
}

type TestData struct {
	Router         *gin.Engine
	DeckService    *DeckService
	FullDeckID     string
	ShuffledDeckID string
	PartialDeckID  string
}

func setupTest(t *testing.T) TestData {
	db, err := repositories.OpenDatabase(&testConfig.DBConfig)
	if err != nil {
		t.Fatalf("Failed to open DB connection: %s", err.Error())
	}
	if err = repositories.AutoMigrateModels(db); err != nil {
		t.Fatalf("Failed to generate db structure: %s", err.Error())
	}

	deckRepo := deck_repo.NewDBDeckRepository(db)
	deckService := NewDeckService(deckRepo)

	router := gin.Default()
	router.Use(errors.ErrorHandler())

	SetupDeckServiceRouting(router.Group("/v1"), deckService)

	fullDeck, err := deckRepo.CreateDeck(&deck_repo.Deck{})
	if err != nil {
		t.Fatalf("Failed to create full deck: %s", err.Error())
	}
	shuffledDeck, err := deckRepo.CreateDeck(&deck_repo.Deck{
		Shuffled: true,
	})
	if err != nil {
		t.Fatalf("Failed to create shuffled deck: %s", err.Error())
	}
	partialDeck, err := deckRepo.CreateDeck(&deck_repo.Deck{
		Cards: []*deck_repo.Card{
			{
				Value: "AC",
			},
			{
				Value: "2C",
			},
			{
				Value: "3C",
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create full deck: %s", err.Error())
	}

	return TestData{
		Router:         router,
		DeckService:    deckService,
		FullDeckID:     fullDeck.ID,
		ShuffledDeckID: shuffledDeck.ID,
		PartialDeckID:  partialDeck.ID,
	}
}

func TestCreateDeck(t *testing.T) {
	td := setupTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/deck", nil)

	td.Router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected response code 201 (Created), but got %d instead.", w.Code)
	}

	resp := &CreateDeckResponse{}
	if err := json.Unmarshal(w.Body.Bytes(), resp); err != nil {
		t.Fatalf("Expected to deserialize the deck, but got error: %s", err.Error())
	}

	if resp.DeckID == "" {
		t.Error("Expected the deck id to be set.")
	}
	if resp.Remaining != 52 {
		t.Error("Expected the deck to have 52 cards remainig.")
	}
	if resp.Shuffled {
		t.Error("Expected the deck not to be shuffled.")
	}
}

func TestCreateDeck_Shuffled(t *testing.T) {
	td := setupTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/deck?shuffled=true", nil)

	td.Router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected response code 201 (Created), but got %d instead.", w.Code)
	}

	resp := &CreateDeckResponse{}
	if err := json.Unmarshal(w.Body.Bytes(), resp); err != nil {
		t.Fatalf("Expected to deserialize the deck, but got error: %s", err.Error())
	}

	if resp.DeckID == "" {
		t.Error("Expected the deck id to be set.")
	}
	if resp.Remaining != 52 {
		t.Error("Expected the deck to have 52 cards remainig.")
	}
	if !resp.Shuffled {
		t.Error("Expected the deck to be shuffled.")
	}
}

func TestCreateDeck_Partial(t *testing.T) {
	td := setupTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/deck?cards=5C,6D,7H,8S", nil)

	td.Router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected response code 201 (Created), but got %d instead.", w.Code)
	}

	resp := &CreateDeckResponse{}
	if err := json.Unmarshal(w.Body.Bytes(), resp); err != nil {
		t.Fatalf("Expected to deserialize the deck, but got error: %s", err.Error())
	}

	if resp.DeckID == "" {
		t.Error("Expected the deck id to be set.")
	}
	if resp.Remaining != 4 {
		t.Error("Expected the deck to have 4 cards remainig.")
	}
}

func TestCreateDeck_InvalidCards(t *testing.T) {
	td := setupTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/deck?cards=2C,JK,TW,4D", nil)

	td.Router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected response code 400 (Bad Request), but got %d instead.", w.Code)
	}

	resp := &errors.ErrorResponse{}
	if err := json.Unmarshal(w.Body.Bytes(), resp); err != nil {
		t.Fatalf("Expected to deserialize the error, but got error: %s", err.Error())
	}
	if resp.Message != "invalid cards values: JK, TW" {
		t.Error("Expected the proper error message for invalid cards.")
	}
}

func TestOpenDeck(t *testing.T) {
	td := setupTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/deck/%s", td.FullDeckID), nil)

	td.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected response code 200 (OK), but got %d instead.", w.Code)
	}

	resp := &OpenDeckResponse{}
	if err := json.Unmarshal(w.Body.Bytes(), resp); err != nil {
		t.Fatalf("Expected to deserialize the deck, but got error: %s", err.Error())
	}

	if resp.DeckID == "" {
		t.Error("Expected the deck id to be set.")
	}
	if resp.Remaining != 52 {
		t.Error("Expected the deck to have 52 cards remainig.")
	}
	if len(resp.Cards) != 52 {
		t.Fatalf("Expected the 52 cards to be shown.")
	}

	if !compare(resp.Cards, strings.Join(deck_repo.NewFullDeck(), ",")) {
		t.Fatalf("Expected to see full deck in order.")
	}
}

func TestOpenDeck_Partial(t *testing.T) {
	td := setupTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/deck/%s", td.PartialDeckID), nil)

	td.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected response code 200 (OK), but got %d instead.", w.Code)
	}

	resp := &OpenDeckResponse{}
	if err := json.Unmarshal(w.Body.Bytes(), resp); err != nil {
		t.Fatalf("Expected to deserialize the deck, but got error: %s", err.Error())
	}

	if resp.DeckID == "" {
		t.Error("Expected the deck id to be set.")
	}
	if resp.Remaining != 3 {
		t.Errorf("Expected the deck to have 3 cards remainig, but has: %d", resp.Remaining)
	}
	if len(resp.Cards) != 3 {
		t.Fatalf("Expected the 3 cards to be shown, but there are actual %d cards.", len(resp.Cards))
	}
}

func TestOpenDeck_NotFound(t *testing.T) {
	td := setupTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/deck/00000000-0000-0000-0000-000000000000", nil)

	td.Router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected response code 404 (Not Found), but got %d instead.", w.Code)
	}

	resp := &errors.ErrorResponse{}
	if err := json.Unmarshal(w.Body.Bytes(), resp); err != nil {
		t.Fatalf("Expected to deserialize the error, but got error: %s", err.Error())
	}
}

func TestDrawCards(t *testing.T) {
	td := setupTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/v1/deck/%s/draw", td.FullDeckID), nil)

	td.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected response code 200 (OK), but got %d instead.", w.Code)
	}

	resp := &DrawCardsResponse{}
	if err := json.Unmarshal(w.Body.Bytes(), resp); err != nil {
		t.Fatalf("Expected to deserialize the deck, but got error: %s", err.Error())
	}

	if len(resp.Cards) != 1 {
		t.Fatalf("Expected to draw exactly one card, but actuall %d are drawn: ", len(resp.Cards))
	}
	if resp.Cards[0].Code != "AC" {
		t.Fatal("Expected to draw the first card.")
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/v1/deck/%s", td.FullDeckID), nil)

	td.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected response code 200 (OK), but got %d instead.", w.Code)
	}

	deck := &OpenDeckResponse{}
	if err := json.Unmarshal(w.Body.Bytes(), deck); err != nil {
		t.Fatalf("Expected to deserialize the deck, but got error: %s", err.Error())
	}

	if deck.Remaining != 51 {
		t.Fatalf("Expected the deck to have 51 cards after drawing 1, but actually has: %d", deck.Remaining)
	}
}

func TestDrawCards_DeckNotFound(t *testing.T) {
	td := setupTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/deck/00000000-0000-0000-0000-000000000000/draw", nil)

	td.Router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected response code 404 (Not Found), but got %d instead.", w.Code)
	}

	resp := &errors.ErrorResponse{}
	if err := json.Unmarshal(w.Body.Bytes(), resp); err != nil {
		t.Fatalf("Expected to deserialize the error, but got error: %s", err.Error())
	}
}

func TestDrawCards_Overdraw(t *testing.T) {
	td := setupTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/v1/deck/%s/draw?count=4", td.PartialDeckID), nil)

	td.Router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected response code 400 (Bad Request), but got %d instead.", w.Code)
	}

	resp := &errors.ErrorResponse{}
	if err := json.Unmarshal(w.Body.Bytes(), resp); err != nil {
		t.Fatalf("Expected to deserialize the error, but got error: %s", err.Error())
	}
	if resp.Message != "not enough cards in deck" {
		t.Errorf("Expected correct message when drawing invalid number of cards, but got instead: '%s'", resp.Message)
	}
}

func compare(deck1 []CardResponse, deck2 string) bool {
	deck1Arr := []string{}
	for _, card := range deck1 {
		deck1Arr = append(deck1Arr, card.Code)
	}
	return strings.Join(deck1Arr, ",") == deck2
}
