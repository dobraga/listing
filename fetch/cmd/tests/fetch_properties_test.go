package main

import (
	"fetch/pkg/domain/property"
	"fetch/pkg/models"
	"fetch/pkg/utils"
	"fmt"
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

	properties, errs := property.FetchProperties(location)
	if errs != nil {
		t.Errorf("%v", errs)
		return
	}
	fmt.Printf("%+v", properties[0])
}

func TestFetchPropertiesStreet(t *testing.T) {
	utils.LoadSettings()
	logrus.SetLevel(logrus.DebugLevel)

	location := models.SearchConfig{
		Local: models.Local{
			City:            "Rio de Janeiro",
			Zone:            "Zona Norte",
			State:           "Rio de Janeiro",
			LocationId:      "BR>Rio de Janeiro>NULL>Rio de Janeiro>Zona Norte>Tijuca",
			Neighborhood:    "Tijuca",
			StateAcronym:    "RJ",
			AddressStreet:   "Rua Barão de Mesquita",
			AddressPointLat: -22.925208,
			AddressPointLon: -43.250265,
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

	properties, errs := property.FetchProperties(location)
	if errs != nil {
		t.Errorf("%v", errs)
		return
	}
	fmt.Printf("%+v", properties[0])
}
