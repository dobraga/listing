package database

import (
	"fetch/property"
	"fetch/station"

	"gorm.io/gorm"
)

func AutoMigrate(DB *gorm.DB) {
	err := DB.AutoMigrate(&property.Property{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&station.Station{})
	if err != nil {
		panic(err)
	}
}
