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
	PartialDeck    string
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
		PartialDeck:    partialDeck.ID,
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

func compare(deck1 []CardResponse, deck2 string) bool {
	deck1Arr := []string{}
	for _, card := range deck1 {
		deck1Arr = append(deck1Arr, card.Code)
	}
	return strings.Join(deck1Arr, ",") == deck2
}
