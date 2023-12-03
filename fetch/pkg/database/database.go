package database

import (
	"database/sql"
	"fetch/pkg/models"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Database struct {
	DB    *gorm.DB
	dbSQL *sql.DB
}

func (db *Database) Ping() error {
	return db.dbSQL.Ping()
}

func (db *Database) Close() {
	db.dbSQL.Close()
}

func (db *Database) AutoMigrate() error {
	err := db.DB.AutoMigrate(&models.Property{})
	if err != nil {
		return err
	}
	return db.DB.AutoMigrate(&models.Station{})
}

func (db *Database) Where(config models.SearchConfig) ([]models.Property, error) {
	var properties []models.Property
	err := db.DB.Where("active = ? AND neighborhood = ? AND city = ? AND state = ? AND business_type = ? AND listing_type = ?",
		true,
		config.Local.Neighborhood,
		config.Local.City,
		config.Local.State,
		config.BusinessType,
		config.ListingType).Find(&properties).Error

	return properties, err
}

func (db *Database) GetLastUpdate(config models.SearchConfig) (time.Time, error) {
	var lastUpdate time.Time

	result := db.DB.Model(&models.Property{}).Where(
		"neighborhood = ? AND city = ? AND state = ? AND business_type = ? AND listing_type = ? AND active",
		config.Local.Neighborhood,
		config.Local.City,
		config.Local.State,
		config.BusinessType,
		config.ListingType).Select("MAX(updated_date)").Row()

	err := result.Scan(&lastUpdate)

	return lastUpdate, err
}

func (db *Database) ResetActive(config models.SearchConfig) error {
	logrus.Debugf(
		"Reset active properties -> neighborhood = '%s' AND city = '%s' AND state = '%s' business_type = '%s' AND listing_type = '%s'",
		config.Local.Neighborhood,
		config.Local.City,
		config.Local.State,
		config.BusinessType,
		config.ListingType)

	tx := db.DB.Model(&models.Property{}).Where(
		"neighborhood = ? AND city = ? AND state = ? AND business_type = ? AND listing_type = ?",
		config.Local.Neighborhood,
		config.Local.City,
		config.Local.State,
		config.BusinessType,
		config.ListingType).Update("Active", false)

	return tx.Error
}

func (db *Database) StoreProperty(config models.SearchConfig, properties []models.Property) error {
	logrus.Infof("Inserting %d properties from to database", len(properties))
	result := db.DB.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(properties, 1)
	return result.Error
}

func (db *Database) StoreStations(stations []models.Station) error {
	return db.DB.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(stations, 50).Error
}
