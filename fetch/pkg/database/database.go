package database

import (
	"database/sql"
	"fetch/pkg/models"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Database struct {
	db    *gorm.DB
	dbSQL *sql.DB
}

func (db *Database) Ping() error {
	return db.dbSQL.Ping()
}

func (db *Database) Close() {
	db.dbSQL.Close()
}

func (db *Database) AutoMigrate() error {
	err := db.db.AutoMigrate(&models.Property{})
	if err != nil {
		return err
	}
	return db.db.AutoMigrate(&models.Station{})
}

func (db *Database) ResetActive(location models.SearchConfig) error {
	logrus.Debugf(
		"Reset active properties -> neighborhood = '%s' AND city = '%s' AND state = '%s' business_type = '%s' AND listing_type = '%s'",
		location.Local.Neighborhood,
		location.Local.City,
		location.Local.State,
		location.BusinessType,
		location.ListingType)

	tx := db.db.Model(&models.Property{}).Where(
		"neighborhood = ? AND city = ? AND state = ? AND business_type = ? AND listing_type = ?",
		location.Local.Neighborhood,
		location.Local.City,
		location.Local.State,
		location.BusinessType,
		location.ListingType).Update("Active", false)

	return tx.Error
}

func (db *Database) StoreProperty(config models.SearchConfig, properties []models.Property) error {
	logrus.Infof("Inserting %d properties from to database", len(properties))
	result := db.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(properties, 1)
	return result.Error
}

func (db *Database) StoreStations(stations []models.Station) error {
	return db.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(stations, 50).Error
}