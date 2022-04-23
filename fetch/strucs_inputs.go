package main

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

var businessTypeValues = map[string]bool{"RENTAL": true, "SALE": true}
var listingTypeValues = map[string]bool{"DEVELOPMENT": true, "USED": true}

type Local struct {
	City         string
	Zone         string
	State        string
	LocationId   string
	Neighborhood string
	StateAcronym string
}

type Location struct {
	Local        Local
	BusinessType string
	ListingType  string
	Origin       string
}

func (l *Location) Validation() error {
	if _, ok := businessTypeValues[l.BusinessType]; !ok {
		return errors.New("business types allowed ['RENTAL', 'SALE']")
	}

	if _, ok := listingTypeValues[l.ListingType]; !ok {
		return errors.New("listing types allowed ['DEVELOPMENT', 'USED']")
	}

	return nil
}

func (l *Location) FinalValidation() []error {
	var errs []error

	err := l.Validation()
	if err != nil {
		errs = append(errs, err)
	}

	if l.Local.City == "" {
		err = fmt.Errorf("need a non empty string into 'City'")
		errs = append(errs, err)
	}

	if l.Local.Zone == "" {
		err = fmt.Errorf("need a non empty string into 'Zone'")
		errs = append(errs, err)
	}

	if l.Local.State == "" {
		err = fmt.Errorf("need a non empty string into 'State'")
		errs = append(errs, err)
	}

	if l.Local.LocationId == "" {
		err = fmt.Errorf("need a non empty string into 'LocationId'")
		errs = append(errs, err)
	}

	if l.Local.Neighborhood == "" {
		err = fmt.Errorf("need a non empty string into 'Neighborhood'")
		errs = append(errs, err)
	}

	if l.Local.StateAcronym == "" {
		err = fmt.Errorf("need a non empty string into 'StateAcronym'")
		errs = append(errs, err)
	}

	// Valida Origem
	configSites := viper.GetStringMapString("sites")
	sites := GetKeys(configSites)

	if !Contains(sites, l.Origin) {
		err = fmt.Errorf("sites need a %v but received '%s'", sites, l.Origin)
		errs = append(errs, err)
	}

	if errs != nil {
		return errs
	}

	return nil
}
