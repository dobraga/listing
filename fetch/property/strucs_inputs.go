package property

import (
	"fmt"
)

var businessTypeValues = map[string]bool{"RENTAL": true, "SALE": true}
var listingTypeValues = map[string]bool{"DEVELOPMENT": true, "USED": true}

type SearchConfig struct {
	Local        Local  `json:"local"`
	BusinessType string `json:"businessType"`
	ListingType  string `json:"listingType"`
	Origin       string `json:"origin"`
}

type Local struct {
	City         string `json:"city"`
	Zone         string `json:"zone"`
	State        string `json:"state"`
	LocationId   string `json:"locationId"`
	Neighborhood string `json:"neighborhood"`
	StateAcronym string `json:"stateAcronym"`
}

func (l *SearchConfig) Validation() []error {
	var errs []error
	var err error

	if _, ok := businessTypeValues[l.BusinessType]; !ok {
		err = fmt.Errorf("business types allowed ['RENTAL', 'SALE']")
		errs = append(errs, err)
	}

	if _, ok := listingTypeValues[l.ListingType]; !ok {
		err = fmt.Errorf("listing types allowed ['DEVELOPMENT', 'USED']")
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
	// configSites := viper.GetStringMapString("sites")
	// sites := utils.GetKeys(configSites)

	// if !utils.Contains(sites, l.Origin) {
	// 	err = fmt.Errorf("sites need a %v but received '%s'", sites, l.Origin)
	// 	errs = append(errs, err)
	// }

	if errs != nil {
		return errs
	}

	return nil
}
