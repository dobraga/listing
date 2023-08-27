package station

import (
	"fetch/pkg/database"

	"github.com/sirupsen/logrus"
)

func CheckExistsExtractCreate(db database.Database) error {
	transportSystems := GetTransportSystems()
	stations, errs := FetchStations(transportSystems)
	logrus.Infof("Extracted %d stations", len(stations))

	if len(errs) > 0 {
		logrus.Errorf("occurs %d errors: %v", len(errs), errs)
	}

	return db.StoreStations(stations)
}
