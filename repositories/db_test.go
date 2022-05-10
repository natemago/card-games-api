package repositories

import (
	"testing"

	"github.com/natemago/card-games-api/config"
)

func TestOpenDatabase(t *testing.T) {
	db, err := OpenDatabase(&config.DBConfig{
		Dialect: "sqlite",
		URL:     "file::memory:?cache=shared",
	})
	if err != nil {
		t.Fatalf("Expected to open a database connection, but got an error instead: %s", err.Error())
	}
	if db == nil {
		t.Error("Expected to get a pointer to gorm.DB")
	}
}

func TestOpenDatabase_UnsupportedDialect(t *testing.T) {
	_, err := OpenDatabase(&config.DBConfig{
		Dialect: "other",
		URL:     "file::memory:?cache=shared",
	})
	if err == nil {
		t.Fatal("Expected to get an unsupported DB type error.")
	}
	if err.Error() != "unsupported DB type: other" {
		t.Error("Expected a valid unsupported DB type error.")
	}
}

func TestAutoMigrateModels(t *testing.T) {
	db, err := OpenDatabase(&config.DBConfig{
		Dialect: "sqlite",
		URL:     "file::memory:?cache=shared",
	})
	if err != nil {
		t.Fatalf("Expected to open a database connection, but got an error instead: %s", err.Error())
	}

	if err = AutoMigrateModels(db); err != nil {
		t.Errorf("Expected to migrate the models, but got an error instead: %s", err.Error())
	}
}
