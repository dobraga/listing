package database

import (
	"fetch/pkg/models"

	"gorm.io/gorm/clause"
)

func StoreStations(stations []models.Station) error {
	db := Connect()
	return db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(stations, 50).Error
}
