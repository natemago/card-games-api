package repositories

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"toggl.com/services/card-games-api/config"
	deck_repo "toggl.com/services/card-games-api/repositories/deck"
)

var MigrationHandlers = []func(db *gorm.DB) error{
	deck_repo.AutoMigrateDeckModels,
}

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
