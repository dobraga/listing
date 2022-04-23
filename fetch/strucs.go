package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Station struct {
	Name     string
	Uf       string `gorm:"primaryKey"`
	Linha    string `gorm:"primaryKey"`
	Station  string `gorm:"primaryKey"`
	Lat      float64
	Lon      float64
	Url      string
	URLLinha string
}

type Property struct {
	Origin        string `gorm:"primaryKey"`
	Url           string `gorm:"primaryKey"`
	Neighborhood  string
	State         string
	City          string
	Zone          string
	BusinessType  string `gorm:"primaryKey"`
	ListingType   string
	Title         string
	UsableArea    int
	Floors        int
	UnitTypes     string
	Bedrooms      int
	Bathrooms     int
	Suites        int
	ParkingSpaces int
	Amenities     string
	Lat           float64
	Lon           float64
	Price         float64
	CondoFee      float64
	CreatedDate   time.Time
	UpdatedDate   time.Time
	Images        string
}

func (p *Property) Unmarshal(bytesData []byte, url, BusinessType string) ([]Property, error) {
	var listNestedProperty []NestedProperty
	var listProperty []Property
	medias := []string{}

	// Bytes to map of interfaces
	data := map[string]interface{}{}
	err := json.Unmarshal(bytesData, &data)
	if err != nil {
		log.Error(fmt.Sprintf("Erro no parse dos dados '%v': %v", bytesData, err))
	}

	// Interface to map and get listings
	data = data["search"].(map[string]interface{})
	data = data["result"].(map[string]interface{})
	listingsPage := data["listings"].([]interface{})

	// Slice of listings to Struct
	jsonBytes, err := json.Marshal(listingsPage)
	if err != nil {
		log.Error(fmt.Sprintf("Erro na transformação para binário '%v': %v", listingsPage, err))
		// os.WriteFile("test.json", jsonBytes, 0666)
	}

	err = json.Unmarshal(jsonBytes, &listNestedProperty)
	if err != nil {
		log.Error(fmt.Sprintf("Erro na transformação para struct: %v", err))
	}

	// Unnest struct
	for _, nestedProperty := range listNestedProperty {
		var property Property

		for _, media := range nestedProperty.Medias {
			medias = append(medias, media.Url)
		}

		media := strings.Join(medias, "|")
		amenities := strings.Join(nestedProperty.Listing.Amenities, "|")

		for _, PricingInfos := range nestedProperty.Listing.PricingInfos {
			if PricingInfos.BusinessType == BusinessType {
				price, err := strconv.ParseFloat(PricingInfos.Price, 32)
				if err != nil {
					log.Error(fmt.Sprintf("Erro na transformação do valor '%s' para float: %v", PricingInfos.Price, err))
				}

				monthlyCondoFee, _ := strconv.ParseFloat(PricingInfos.MonthlyCondoFee, 32)

				strUsableArea := GetFirst(nestedProperty.Listing.UsableAreas, property.Url, "UsableAreas")
				usableArea, _ := strconv.Atoi(strUsableArea)

				property.Origin = url
				property.Url = nestedProperty.Link.Href
				property.Neighborhood = nestedProperty.Listing.Address.Neighborhood
				property.State = nestedProperty.Listing.Address.State
				property.City = nestedProperty.Listing.Address.City
				property.Zone = nestedProperty.Listing.Address.Zone
				property.BusinessType = BusinessType
				property.ListingType = nestedProperty.Listing.ListingType
				property.Title = nestedProperty.Listing.Title
				property.UsableArea = usableArea
				property.Floors = GetFirst(nestedProperty.Listing.Floors, property.Url, "Floors")
				property.UnitTypes = GetFirst(nestedProperty.Listing.UnitTypes, property.Url, "UnitTypes")
				property.Bedrooms = GetFirst(nestedProperty.Listing.Bedrooms, property.Url, "Bedrooms")
				property.Bathrooms = GetFirst(nestedProperty.Listing.Bathrooms, property.Url, "Bathrooms")
				property.Suites = GetFirst(nestedProperty.Listing.Suites, property.Url, "Suites")
				property.ParkingSpaces = GetFirst(nestedProperty.Listing.ParkingSpaces, property.Url, "ParkingSpaces")
				property.Amenities = amenities
				property.Lat = float64(nestedProperty.Listing.Address.Point.Lat)
				property.Lon = float64(nestedProperty.Listing.Address.Point.Lon)
				property.Price = price
				property.CondoFee = monthlyCondoFee
				property.CreatedDate = nestedProperty.Listing.CreatedAt
				property.UpdatedDate = nestedProperty.Listing.UpdatedAt
				property.Images = media

				listProperty = append(listProperty, property)
			}
		}

	}

	return listProperty, nil
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
	StreetNumber string  `json:"streetNumber"`
	ZipCode      string  `json:"zipCode"`
	Zone         string  `json:"zone"`
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
