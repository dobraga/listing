package main

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ListLocations(local string, origin string) ([]Location, error) {
	var final_locations []Location

	sites := viper.Get("sites")
	if sites == nil {
		return nil, errors.New("variável de sites vazia")
	}

	site_info := sites.(map[string]interface{})[origin].(map[string]interface{})
	portal := site_info["portal"].(string)

	query := map[string]interface{}{
		"q": local, "fields": "neighborhood", "portal": portal, "size": "6",
	}

	bytes_data := MakeRequest(true, origin, query)

	data := map[string]interface{}{}
	err := json.Unmarshal(bytes_data, &data)
	if err != nil {
		log.Error(fmt.Sprintf("Erro ao listar localizações %s: %v", local, err))
	}

	// Interface to map and get listings
	data = data["neighborhood"].(map[string]interface{})
	data = data["result"].(map[string]interface{})
	locations := data["locations"].([]interface{})

	for _, location := range locations {
		address := location.(map[string]interface{})["address"].(map[string]interface{})

		s_location := Location{
			City:         address["city"].(string),
			Zone:         address["zone"].(string),
			State:        address["state"].(string),
			LocationId:   address["locationId"].(string),
			Neighborhood: address["neighborhood"].(string),
			StateAcronym: address["stateAcronym"].(string),
		}

		final_locations = append(final_locations, s_location)
	}

	return final_locations, nil
}
