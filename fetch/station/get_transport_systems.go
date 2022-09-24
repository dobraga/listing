package station

import (
	"github.com/spf13/viper"
)

func GetTransportSystems() []TransportSystem {
	var transportSystems []TransportSystem
	mapSites := viper.GetStringMap("metro_trem")

	for uf, urls := range mapSites {
		listUrls := urls.([]interface{})

		for _, url := range listUrls {
			transportSystems = append(transportSystems, TransportSystem{URL: url.(string), Uf: uf})
		}
	}

	return transportSystems
}
