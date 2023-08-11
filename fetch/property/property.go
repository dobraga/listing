package property

import (
	"encoding/json"
	"fetch/models"
	"fetch/utils"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func UnmarshalProperty(data map[string]interface{}, l models.SearchConfig) ([]models.Property, error) {
	var listNestedProperty []models.NestedProperty
	var listProperty []models.Property
	medias := []string{}

	// Interface to map and get listings
	if utils.Contains(utils.GetKeys(data), "nearby") {
		data = data["nearby"].(map[string]interface{})
	}
	data = data["search"].(map[string]interface{})
	data = data["result"].(map[string]interface{})
	listingsPage := data["listings"].([]interface{})

	// Slice of listings to Struct
	jsonBytes, err := json.Marshal(listingsPage)
	// os.WriteFile("test.json", jsonBytes, 0666)
	if err != nil {
		err = fmt.Errorf("erro na transformação para binário '%v': %v", listingsPage, err)
		log.Error(err)
		return listProperty, err
	}

	err = json.Unmarshal(jsonBytes, &listNestedProperty)
	if err != nil {
		err = fmt.Errorf("erro na transformação para struct: %v", err)
		log.Error(err)
		return listProperty, err
	}

	// Unnest struct
	for _, nestedProperty := range listNestedProperty {
		log.Debugf("Processing: '%s'", nestedProperty.Link.Href)

		var property models.Property

		for _, media := range nestedProperty.Medias {
			medias = append(medias, media.Url)
		}

		media := strings.Join(medias, "|")
		amenities := strings.Join(nestedProperty.Listing.Amenities, "|")

		for _, PricingInfos := range nestedProperty.Listing.PricingInfos {
			if PricingInfos.BusinessType == l.BusinessType {
				price, err := strconv.ParseFloat(PricingInfos.Price, 32)
				if err != nil {
					err = fmt.Errorf("erro na transformação do valor '%s' para float: %v", PricingInfos.Price, err)
					log.Error(err)
					return listProperty, err
				}

				monthlyCondoFee, _ := strconv.ParseFloat(PricingInfos.MonthlyCondoFee, 32)

				strUsableArea := utils.GetFirst(nestedProperty.Listing.UsableAreas, property.Url, "UsableAreas")
				usableArea, _ := strconv.Atoi(strUsableArea)

				property.Origin = l.Origin
				property.Url = nestedProperty.Link.Href
				property.Neighborhood = nestedProperty.Listing.Address.Neighborhood
				property.State = nestedProperty.Listing.Address.State
				property.StateAcronym = nestedProperty.Listing.Address.StateAcronym
				property.City = nestedProperty.Listing.Address.City
				property.Zone = nestedProperty.Listing.Address.Zone
				property.Street = nestedProperty.Listing.Address.Street
				property.StreetNumber = nestedProperty.Listing.Address.StreetNumber
				property.BusinessType = l.BusinessType
				property.ListingType = nestedProperty.Listing.ListingType
				property.Title = nestedProperty.Listing.Title
				property.UsableArea = usableArea
				property.Floors = utils.GetFirst(nestedProperty.Listing.Floors, property.Url, "Floors")
				property.UnitTypes = utils.GetFirst(nestedProperty.Listing.UnitTypes, property.Url, "UnitTypes")
				property.Bedrooms = utils.GetFirst(nestedProperty.Listing.Bedrooms, property.Url, "Bedrooms")
				property.Bathrooms = utils.GetFirst(nestedProperty.Listing.Bathrooms, property.Url, "Bathrooms")
				property.Suites = utils.GetFirst(nestedProperty.Listing.Suites, property.Url, "Suites")
				property.ParkingSpaces = utils.GetFirst(nestedProperty.Listing.ParkingSpaces, property.Url, "ParkingSpaces")
				property.Amenities = amenities
				property.Lat = float64(nestedProperty.Listing.Address.Point.Lat)
				property.Lon = float64(nestedProperty.Listing.Address.Point.Lon)
				property.Price = price
				property.CondoFee = monthlyCondoFee
				property.CreatedDate = nestedProperty.Listing.CreatedAt
				property.UpdatedDate = nestedProperty.Listing.UpdatedAt
				if !l.DropImages {
					property.Images = media
				}
				property.Active = true

				listProperty = append(listProperty, property)
			}
		}

		log.Debugf("Processed: '%s'", nestedProperty.Link.Href)
	}

	return listProperty, nil
}
