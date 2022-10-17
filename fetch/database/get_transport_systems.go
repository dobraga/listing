package database

import (
	"fetch/models"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GetTransportSystems() []models.TransportSystem {
	var allTransportSystems []models.TransportSystem
	var transportSystems []models.TransportSystem

	mapSites := viper.GetStringMap("metro_trem")
	db := Connect()

	for uf, urls := range mapSites {
		listUrls := urls.([]interface{})

		for _, url := range listUrls {
			allTransportSystems = append(
				allTransportSystems,
				models.TransportSystem{URL: url.(string), Uf: uf})
		}
	}

	for _, system := range allTransportSystems {
		var station models.Station

		db.Where("transport_system_url = ?", system.URL).First(&station)
		if station.TransportSystemURL != system.URL {
			logrus.Infof("Extract data from '%s'", system.URL)
			transportSystems = append(transportSystems, system)
		} else {
			logrus.Debugf("Already extracted data from '%s'", system.URL)
		}
	}

	return transportSystems
}
