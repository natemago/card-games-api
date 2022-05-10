package app

import (
	"github.com/natemago/card-games-api/config"
	"github.com/natemago/card-games-api/repositories"
	deck_repo "github.com/natemago/card-games-api/repositories/deck"
	"github.com/natemago/card-games-api/rest"
	deck_svcs "github.com/natemago/card-games-api/rest/deck"
)

func RunApp(conf *config.Config) error {
	// Build and connect to database
	db, err := repositories.OpenDatabase(&conf.DBConfig)
	if err != nil {
		return err
	}

	// Do model migration
	if err := repositories.AutoMigrateModels(db); err != nil {
		return err
	}

	// Build the repositories
	deckRepository := deck_repo.NewDBDeckRepository(db)

	// Buld the services
	deckService := deck_svcs.NewDeckService(deckRepository)

	// Finally run the API
	return rest.RunAPI(&conf.APIConfig, *deckService)
}
