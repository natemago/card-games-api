package app

import (
	"toggl.com/services/card-games-api/config"
	"toggl.com/services/card-games-api/repositories"
	deck_repo "toggl.com/services/card-games-api/repositories/deck"
	"toggl.com/services/card-games-api/rest"
	deck_svcs "toggl.com/services/card-games-api/rest/deck"
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
