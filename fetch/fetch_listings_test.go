package main

import (
	"testing"
)

func TestFetchListing(t *testing.T) {
	LoadSettings()

	// location := Location{}
	// errs := location.FinalValidation()

	// t.Errorf("%v", errs)

	location := Location{
		Local: Local{
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

	errs := location.Validation()
	if errs != nil {
		t.Errorf("%v", errs)
	}

	db := Connect()

	_, errs = FetchListings(db, location)
	if errs != nil {
		t.Errorf("%v", errs)
	}
}
