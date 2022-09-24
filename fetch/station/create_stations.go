package station

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CheckExistsExtractCreate(db *gorm.DB) error {
	var transportSystems []TransportSystem

	allTransportSystems := GetTransportSystems()

	for _, system := range allTransportSystems {
		var station Station

		db.Where("transport_system_url = ?", system.URL).First(&station)
		if station.TransportSystemURL != system.URL {
			logrus.Infof("Extract data from '%s'", system.URL)
			transportSystems = append(transportSystems, system)
		} else {
			logrus.Debugf("Already extracted data from '%s'", system.URL)
		}
	}

	stations, errs := FetchStations(transportSystems)
	logrus.Infof("Extracted %d stations", len(stations))

	if len(errs) > 0 {
		logrus.Errorf("occurs %d errors: %v", len(errs), errs)
	}

	if len(stations) > 0 {
		result := db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(stations, 500)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
