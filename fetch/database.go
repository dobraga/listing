package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	var db *gorm.DB
	var err error

	if viper.GetString("ENV") == "DEVELOPMENT" {
		db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
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
		log.Fatal(err)
	}

	db.AutoMigrate(&Station{})
	db.AutoMigrate(&Property{})

	return db
}
