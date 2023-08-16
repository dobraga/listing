package models

import "time"

type Property struct {
	Origin            string    `gorm:"primaryKey" json:"-"`
	Url               string    `gorm:"primaryKey" json:"url"`
	Neighborhood      string    `json:"neighborhood"`
	State             string    `json:"-"`
	StateAcronym      string    `json:"-"`
	City              string    `json:"-"`
	Zone              string    `json:"-"`
	Street            string    `json:"-"`
	StreetNumber      string    `json:"-"`
	BusinessType      string    `gorm:"primaryKey" json:"-"`
	ListingType       string    `json:"-"`
	Title             string    `json:"-"`
	UsableArea        int       `json:"usable_area"`
	Floors            int       `json:"-"`
	UnitTypes         string    `json:"unit_types"`
	Bedrooms          int       `json:"bedrooms"`
	Bathrooms         int       `json:"bathrooms"`
	Suites            int       `json:"-"`
	ParkingSpaces     int       `json:"parking_spaces"`
	Amenities         string    `json:"-"`
	Lat               float64   `json:"lat"`
	Lon               float64   `json:"lon"`
	Price             float64   `json:"-"`
	CondoFee          float64   `json:"-"`
	PredictTotalPrice float64   `json:"-"`
	CreatedDate       time.Time `json:"-"`
	UpdatedDate       time.Time `json:"-"`
	Images            string    `json:"-"`
	Active            bool      `json:"-"`
}

type NestedProperty struct {
	Listing Listing `json:"listing"`
	Medias  []Media `json:"medias"`
	Link    Link    `json:"link"`
}

type Media struct {
	Url string `json:"url"`
}

type Link struct {
	Href string `json:"href"`
}

type Listing struct {
	Id              string    `json:"id"`
	Title           string    `json:"title"`
	UsableAreas     []string  `json:"usableAreas"`
	Address         Address   `json:"address"`
	Amenities       []string  `json:"amenities"`
	Bathrooms       []int     `json:"bathrooms"`
	Bedrooms        []int     `json:"bedrooms"`
	Suites          []int     `json:"suites"`
	Description     string    `json:"description"`
	Floors          []int     `json:"floors"`
	ListingType     string    `json:"listingType"`
	UsageTypes      []string  `json:"usageTypes"`
	ParkingSpaces   []int     `json:"parkingSpaces"`
	Portal          string    `json:"portal"`
	PricingInfos    []Pricing `json:"pricingInfos"`
	PropertyType    string    `json:"propertyType"`
	PublicationType string    `json:"publicationType"`
	UnitTypes       []string  `json:"unitTypes"`
	UnitsOnTheFloor int       `json:"unitsOnTheFloor"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type Address struct {
	City         string  `json:"city"`
	Complement   string  `json:"complement"`
	LocationId   string  `json:"locationId"`
	Neighborhood string  `json:"neighborhood"`
	Point        LatLong `json:"point"`
	State        string  `json:"state"`
	StateAcronym string  `json:"stateAcronym"`
	ZipCode      string  `json:"zipCode"`
	Zone         string  `json:"zone"`
	Street       string  `json:"street"`
	StreetNumber string  `json:"streetNumber"`
}

type LatLong struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type Pricing struct {
	BusinessType    string `json:"businessType"`
	MonthlyCondoFee string `json:"monthlyCondoFee"`
	Price           string `json:"price"`
	YearlyIptu      string `json:"yearlyIptu"`
}
