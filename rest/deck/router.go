package deck

import "github.com/gin-gonic/gin"

// SetupDeckServiceRouting sets up the routing for DeckService with gin router.
func SetupDeckServiceRouting(group *gin.RouterGroup, deckService *DeckService) {
	group.POST("/deck", deckService.CreateDeck)
	group.GET("/deck/:deckId", deckService.OpenDeck)
	group.POST("/deck/:deckId/draw", deckService.DrawCards)
}
