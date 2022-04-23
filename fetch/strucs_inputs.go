package main

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

var business_type_values = map[string]bool{"RENTAL": true, "SALE": true}
var listing_type_values = map[string]bool{"DEVELOPMENT": true, "USED": true}

type Location struct {
	City         string
	Zone         string
	State        string
	LocationId   string
	Neighborhood string
	StateAcronym string
	BusinessType string
	ListingType  string
	Origin       string
}

func (l *Location) Validation() error {
	if _, ok := business_type_values[l.BusinessType]; !ok {
		return errors.New("business types allowed ['RENTAL', 'SALE']")
	}

	if _, ok := listing_type_values[l.ListingType]; !ok {
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

	if l.City == "" {
		err = fmt.Errorf("need a non empty string into 'City'")
		errs = append(errs, err)
	}

	if l.Zone == "" {
		err = fmt.Errorf("need a non empty string into 'Zone'")
		errs = append(errs, err)
	}

	if l.State == "" {
		err = fmt.Errorf("need a non empty string into 'State'")
		errs = append(errs, err)
	}

	if l.LocationId == "" {
		err = fmt.Errorf("need a non empty string into 'LocationId'")
		errs = append(errs, err)
	}

	if l.Neighborhood == "" {
		err = fmt.Errorf("need a non empty string into 'Neighborhood'")
		errs = append(errs, err)
	}

	if l.StateAcronym == "" {
		err = fmt.Errorf("need a non empty string into 'StateAcronym'")
		errs = append(errs, err)
	}

	// Valida Origem
	config_sites := viper.GetStringMapString("sites")
	sites := GetKeys(config_sites)

	if !Contains(sites, l.Origin) {
		err = fmt.Errorf("sites need a %v but received '%s'", sites, l.Origin)
		errs = append(errs, err)
	}

	if errs != nil {
		return errs
	}

	return nil
}
