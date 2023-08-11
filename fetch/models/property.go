package models

import "time"

type Property struct {
	Origin            string `gorm:"primaryKey"`
	Url               string `gorm:"primaryKey"`
	Neighborhood      string
	State             string
	StateAcronym      string
	City              string
	Zone              string
	Street            string
	StreetNumber      string
	BusinessType      string `gorm:"primaryKey"`
	ListingType       string
	Title             string
	UsableArea        int
	Floors            int
	UnitTypes         string
	Bedrooms          int
	Bathrooms         int
	Suites            int
	ParkingSpaces     int
	Amenities         string
	Lat               float64
	Lon               float64
	Price             float64
	CondoFee          float64
	PredictTotalPrice float64
	CreatedDate       time.Time
	UpdatedDate       time.Time
	Images            string
	Active            bool
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
