package database

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() (Database, error) {
	var DB Database
	var db *gorm.DB
	var err error

	if viper.GetString("DEVELOPMENT_DATABASE") != "postgres" {
		dsn := sqlite.Open("../test.db")
		db, err = gorm.Open(dsn, &gorm.Config{})
	} else {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
			viper.GetString("POSTGRES_HOST"),
			viper.GetString("POSTGRES_USER"),
			viper.GetString("POSTGRES_PASSWORD"),
			viper.GetString("POSTGRES_DB"),
			viper.GetString("POSTGRES_PORT"),
		)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		return DB, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return DB, err
	}

	database := Database{db, sqlDB}

	return database, database.Ping()
}
