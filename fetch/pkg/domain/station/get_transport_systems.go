package station

import (
	"fetch/pkg/models"

	"github.com/spf13/viper"
)

func GetTransportSystems() []models.TransportSystem {
	var allTransportSystems []models.TransportSystem
	var transportSystems []models.TransportSystem

	mapSites := viper.GetStringMap("metro_trem")

	for uf, urls := range mapSites {
		listUrls := urls.([]interface{})

		for _, url := range listUrls {
			allTransportSystems = append(
				allTransportSystems,
				models.TransportSystem{URL: url.(string), Uf: uf})
		}
	}

	transportSystems = append(transportSystems, allTransportSystems...)

	return transportSystems
}
