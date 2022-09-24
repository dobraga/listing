package main

import (
	"fetch/property"
	"fetch/utils"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestListLocations(t *testing.T) {
	utils.LoadSettings()
	locations, err := property.ListLocations("Tijuca", "vivareal")

	if err != nil {
		t.Error(err)
	}

	if len(locations) != 4 {
		t.Errorf("Expected 4 elements and find %d: %v", len(locations), locations)
	}

	logrus.Infof("Find '%+v' in locations", locations)
}
