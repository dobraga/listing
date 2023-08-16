package database

import (
	"fetch/pkg/models"
)

func AutoMigrate() {
	DB := Connect()

	err := DB.AutoMigrate(&models.Property{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&models.Station{})
	if err != nil {
		panic(err)
	}
}
