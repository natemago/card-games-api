package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/natemago/card-games-api/config"
	"github.com/natemago/card-games-api/errors"
	deck_api "github.com/natemago/card-games-api/rest/deck"
)

// SetupRouting sets up the routing for the whole API.
func SetupRouting(router *gin.Engine, deckService deck_api.DeckService) {
	v1group := router.Group("/v1")
	deck_api.SetupDeckServiceRouting(v1group, &deckService)
}

// SetupAPI sets up the gin router for the whole API based on the provided API Configuration.
func SetupAPI(conf *config.APIConfig) *gin.Engine {
	router := gin.Default()

	router.Use(errors.ErrorHandler())

	return router
}

// RunAPI sets up the API, then sets up routing with gin and finally binds and runs the gin router binding to host and port.
// This will basically set up the whole API, then run the HTTP server to accept connections.
// Blocks until the server is shut down.
func RunAPI(conf *config.APIConfig, deckService deck_api.DeckService) error {
	router := SetupAPI(conf)

	SetupRouting(router, deckService)

	bindAddress := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	return router.Run(bindAddress)
}
