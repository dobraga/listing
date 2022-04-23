package main

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ListLocations(local string, origin string) ([]Location, error) {
	var finalLocations []Location

	sites := viper.Get("sites")
	if sites == nil {
		return nil, errors.New("variável de sites vazia")
	}

	siteInfo := sites.(map[string]interface{})[origin].(map[string]interface{})
	portal := siteInfo["portal"].(string)

	query := map[string]interface{}{
		"q": local, "fields": "neighborhood", "portal": portal, "size": "6",
	}

	bytesData := MakeRequest(true, origin, query)

	data := map[string]interface{}{}
	err := json.Unmarshal(bytesData, &data)
	if err != nil {
		log.Error(fmt.Sprintf("Erro ao listar localizações %s: %v", local, err))
	}

	// Interface to map and get listings
	data = data["neighborhood"].(map[string]interface{})
	data = data["result"].(map[string]interface{})
	locations := data["locations"].([]interface{})

	for _, location := range locations {
		address := location.(map[string]interface{})["address"].(map[string]interface{})

		sLocation := Location{
			City:         address["city"].(string),
			Zone:         address["zone"].(string),
			State:        address["state"].(string),
			LocationId:   address["locationId"].(string),
			Neighborhood: address["neighborhood"].(string),
			StateAcronym: address["stateAcronym"].(string),
		}

		finalLocations = append(finalLocations, sLocation)
	}

	return finalLocations, nil
}
