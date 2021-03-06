package main

import (
	"testing"
)

func TestListLocations(t *testing.T) {
	LoadSettings()
	locations, err := ListLocations("Tijuca", "vivareal")

	if err != nil {
		t.Error(err)
	}

	if len(locations) != 4 {
		t.Errorf("Expected 4 elements and find %d: %v", len(locations), locations)
	}
}
