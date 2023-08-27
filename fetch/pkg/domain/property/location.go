package property

import (
	"errors"
	"fetch/pkg/models"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ListLocations(local, type_location, origin string) ([]models.Local, int, error) {
	logrus.Infof("Searching '%s' locations for type '%s'", local, type_location)
	mapLocations := map[string]models.Local{}
	finalLocations := []models.Local{}

	sites := viper.Get("sites")
	if sites == nil {
		return nil, 400, errors.New("variável de sites vazia")
	}

	siteInfo := sites.(map[string]interface{})[origin].(map[string]interface{})
	portal := siteInfo["portal"].(string)

	query := map[string]interface{}{
		"q": local, "portal": portal, "size": "6",
		"fields":        "neighborhood,city,street",
		"includeFields": "address.city,address.zone,address.state,address.neighborhood,address.stateAcronym,address.street,address.locationId,address.point",
	}

	data, status_code, err := MakeRequest(true, origin, query)
	if err != nil {
		return nil, status_code, err
	}

	// Interface to map and get listings
	data = data["street"].(map[string]interface{})
	data = data["result"].(map[string]interface{})
	locations := data["locations"].([]interface{})

	for _, location := range locations {
		address := location.(map[string]interface{})["address"].(map[string]interface{})
		locationId := address["locationId"].(string)
		point := address["point"].(map[string]interface{})
		sLocation := models.Local{
			City:         address["city"].(string),
			Zone:         address["zone"].(string),
			State:        address["state"].(string),
			LocationId:   locationId,
			Neighborhood: address["neighborhood"].(string),
			StateAcronym: address["stateAcronym"].(string),
		}

		if type_location == "street" {
			sLocation.AddressStreet = address["street"].(string)
			sLocation.AddressPointLat = point["lat"].(float64)
			sLocation.AddressPointLon = point["lon"].(float64)
		}

		mapLocations[locationId] = sLocation
	}

	for _, sLocation := range mapLocations {
		finalLocations = append(finalLocations, sLocation)
	}

	return finalLocations, 200, nil
}
