package database

import (
	"fetch/pkg/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

func StoreProperty(config models.SearchConfig, properties []models.Property) error {
	db := Connect()

	logrus.Infof("Inserting %d properties from '%s' to database", len(properties), config.Origin)
	result := db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(properties, 1)
	return result.Error
}
