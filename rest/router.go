package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/natemago/card-games-api/config"
	deck_api "github.com/natemago/card-games-api/rest/deck"
)

func SetupRouting(router *gin.Engine, deckService deck_api.DeckService) {
	v1group := router.Group("/v1")
	deck_api.SetupDeckServiceRouting(v1group, &deckService)
}

func SetupAPI(conf *config.APIConfig) *gin.Engine {
	router := gin.Default()

	router.Use(ErrorHandler())

	return router
}

func RunAPI(conf *config.APIConfig, deckService deck_api.DeckService) error {
	router := SetupAPI(conf)

	SetupRouting(router, deckService)

	bindAddress := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	return router.Run(bindAddress)
}
