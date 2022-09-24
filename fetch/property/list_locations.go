package property

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

func ListLocations(local string, origin string) ([]Local, error) {
	mapLocations := map[string]Local{}
	finalLocations := []Local{}

	sites := viper.Get("sites")
	if sites == nil {
		return nil, errors.New("variável de sites vazia")
	}

	siteInfo := sites.(map[string]interface{})[origin].(map[string]interface{})
	portal := siteInfo["portal"].(string)

	query := map[string]interface{}{
		"q": local, "fields": "neighborhood", "portal": portal, "size": "6",
	}

	bytesData, err := MakeRequest(true, origin, query)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(bytesData, &data)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar localizações %s: %v", local, err)
	}

	// Interface to map and get listings
	data = data["neighborhood"].(map[string]interface{})
	data = data["result"].(map[string]interface{})
	locations := data["locations"].([]interface{})

	for _, location := range locations {
		address := location.(map[string]interface{})["address"].(map[string]interface{})
		locationId := address["locationId"].(string)

		sLocation := Local{
			City:         address["city"].(string),
			Zone:         address["zone"].(string),
			State:        address["state"].(string),
			LocationId:   locationId,
			Neighborhood: address["neighborhood"].(string),
			StateAcronym: address["stateAcronym"].(string),
		}

		mapLocations[locationId] = sLocation
	}

	for _, sLocation := range mapLocations {
		finalLocations = append(finalLocations, sLocation)
	}

	return finalLocations, nil
}
