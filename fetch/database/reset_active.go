package database

import (
	"fetch/models"

	"github.com/sirupsen/logrus"
)

func ResetActive(location models.SearchConfig) error {
	db := Connect()

	logrus.Debugf(
		"Reset active properties -> neighborhood = '%s' AND city = '%s' AND state = '%s' business_type = '%s' AND listing_type = '%s'",
		location.Local.Neighborhood,
		location.Local.City,
		location.Local.State,
		location.BusinessType,
		location.ListingType)

	tx := db.Model(&models.Property{}).Where(
		"neighborhood = ? AND city = ? AND state = ? AND business_type = ? AND listing_type = ?",
		location.Local.Neighborhood,
		location.Local.City,
		location.Local.State,
		location.BusinessType,
		location.ListingType).Update("Active", false)

	return tx.Error
}
