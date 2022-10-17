package station

import (
	"fetch/database"

	"github.com/sirupsen/logrus"
)

func CheckExistsExtractCreate() error {
	transportSystems := database.GetTransportSystems()
	stations, errs := FetchStations(transportSystems)
	logrus.Infof("Extracted %d stations", len(stations))

	if len(errs) > 0 {
		logrus.Errorf("occurs %d errors: %v", len(errs), errs)
	}

	return database.StoreStations(stations)
}
