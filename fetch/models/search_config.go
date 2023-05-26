package models

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

var businessTypeValues = map[string]bool{"RENTAL": true, "SALE": true}
var listingTypeValues = map[string]bool{"DEVELOPMENT": true, "USED": true}

type SearchConfig struct {
	Local        Local  `json:"local"`
	BusinessType string `json:"businessType"`
	ListingType  string `json:"listingType"`
	Origin       string `json:"origin"`
}

func (s *SearchConfig) ExtractFromContext(c *gin.Context) {
	addressPointLat, _ := strconv.ParseFloat(c.Param("addressPointLat"), 64)
	addressPointLon, _ := strconv.ParseFloat(c.Param("addressPointLon"), 64)

	s.Local = Local{
		City:            c.Param("city"),
		Zone:            c.Param("zone"),
		State:           c.Param("state"),
		LocationId:      c.Param("locationId"),
		Neighborhood:    c.Param("neighborhood"),
		StateAcronym:    c.Param("stateAcronym"),
		AddressStreet:   c.Param("addressStreet"),
		AddressPointLat: addressPointLat,
		AddressPointLon: addressPointLon,
	}
	s.BusinessType = c.Param("business_type")
	s.ListingType = c.Param("listing_type")
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

type Local struct {
	City            string  `json:"city"`
	Zone            string  `json:"zone"`
	State           string  `json:"state"`
	LocationId      string  `json:"locationId"`
	Neighborhood    string  `json:"neighborhood"`
	StateAcronym    string  `json:"stateAcronym"`
	AddressStreet   string  `json:"addressStreet,omitempty"`
	AddressPointLat float64 `json:"addressPointLat,omitempty"`
	AddressPointLon float64 `json:"addressPointLon,omitempty"`
}
