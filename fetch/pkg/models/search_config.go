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
	DropImages   bool   `json:"dropImages"`
}

func (s *SearchConfig) ExtractFromContext(c *gin.Context) {
	dropImages, _ := strconv.ParseBool(c.DefaultQuery("dropImages", "false"))
	s.Local = Local{
		City:         c.Query("city"),
		Zone:         c.Query("zone"),
		State:        c.Query("state"),
		LocationId:   c.Query("locationId"),
		Neighborhood: c.Query("neighborhood"),
		StateAcronym: c.Query("stateAcronym"),
	}
	s.BusinessType = c.Query("business_type")
	s.ListingType = c.Query("listing_type")
	s.DropImages = dropImages

	addressStreet := c.Query("addressStreet")
	if addressStreet != "" {
		s.Local.AddressStreet = addressStreet
		addressPointLat, _ := strconv.ParseFloat(c.Query("addressPointLat"), 64)
		addressPointLon, _ := strconv.ParseFloat(c.Query("addressPointLon"), 64)
		s.Local.AddressPointLat = addressPointLat
		s.Local.AddressPointLon = addressPointLon
	}
}

func (l *SearchConfig) Validation() []error {
	var errs []error
	var err error

	if _, ok := businessTypeValues[l.BusinessType]; !ok {
		err = fmt.Errorf("business types allowed [business_type='RENTAL' or business_type='SALE']")
		errs = append(errs, err)
	}

	if _, ok := listingTypeValues[l.ListingType]; !ok {
		err = fmt.Errorf("listing types allowed [listing_type='DEVELOPMENT' or listing_type='USED']")
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

	return errs
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
