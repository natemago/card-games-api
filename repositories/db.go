package repositories

import (
	"fmt"

	"github.com/natemago/card-games-api/config"
	deck_repo "github.com/natemago/card-games-api/repositories/deck"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MigrationHandlers is a list of functions that perform a database migrations for multiple models.
var MigrationHandlers = []func(db *gorm.DB) error{
	deck_repo.AutoMigrateDeckModels,
}

// OpenDatabase creates a new connection to the database based on the supplied database configuration config.DBConfig.
func OpenDatabase(config *config.DBConfig) (db *gorm.DB, err error) {
	dialect, err := configureDialect(config)
	if err != nil {
		return nil, err
	}
	db, err = gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func configureDialect(config *config.DBConfig) (gorm.Dialector, error) {
	switch config.Dialect {
	case "postgres":
		return postgres.Open(config.URL), nil
	case "sqlite":
		return sqlite.Open(config.URL), nil
	default:
		return nil, fmt.Errorf("unsupported DB type: %s", config.Dialect)
	}
}

// AutoMigrateModels executes the MigrationHandlers within a single transaction to auto-migrate the database models.
func AutoMigrateModels(db *gorm.DB) error {
	tx := db.Begin()

	for _, autoMigrateHandler := range MigrationHandlers {
		if err := autoMigrateHandler(tx); err != nil {
			result := tx.Rollback()
			if result.Error != nil {
				return fmt.Errorf("failed to rollback migration transaction: %s; original error: %s", result.Error.Error(), err.Error())
			}
			return err
		}
	}
	result := tx.Commit()
	if result.Error != nil {
		return result.Error
	}
	return nil
}
