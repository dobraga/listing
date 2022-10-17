package main

import (
	"fetch/models"
	"fetch/property"
	"fetch/utils"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestFetchProperties(t *testing.T) {
	utils.LoadSettings()
	logrus.SetLevel(logrus.DebugLevel)

	location := models.SearchConfig{
		Local: models.Local{
			City:         "Rio de Janeiro",
			Zone:         "Zona Norte",
			State:        "Rio de Janeiro",
			LocationId:   "BR>Rio de Janeiro>NULL>Rio de Janeiro>Zona Norte>Tijuca",
			Neighborhood: "Tijuca",
			StateAcronym: "RJ",
		},
		BusinessType: "RENTAL",
		ListingType:  "USED",
		Origin:       "vivareal",
	}
	err := location.Validation()
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	errs := property.FetchProperties(location, 0)
	if errs != nil {
		t.Errorf("%v", errs)
		return
	}
}
