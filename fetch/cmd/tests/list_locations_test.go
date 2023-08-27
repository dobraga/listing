package main

import (
	"fetch/pkg/domain/property"
	"fetch/pkg/utils"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestListLocations(t *testing.T) {
	utils.LoadSettings()
	locations, _, err := property.ListLocations("Tijuca", "neighborhood", "vivareal")

	if err != nil {
		t.Error(err)
	}

	if len(locations) != 2 {
		t.Errorf("Expected 2 elements and find %d: %v", len(locations), locations)
	}

	logrus.Infof("Find '%+v' in locations", locations)
}

func TestListStreet(t *testing.T) {
	utils.LoadSettings()
	locations, _, err := property.ListLocations("Rua Barão", "street", "vivareal")

	if err != nil {
		t.Error(err)
	}

	if len(locations) != 5 {
		t.Errorf("Expected 5 element and find %d: %v", len(locations), locations)
	}

	logrus.Infof("Find '%+v' in locations", locations)
}
