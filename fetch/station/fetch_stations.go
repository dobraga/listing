package station

import (
	"fetch/models"

	"github.com/sirupsen/logrus"
)

func FetchStations(transports_systems []models.TransportSystem) ([]models.Station, []error) {
	var stations []models.Station
	var lines []models.Line
	var errs []error

	logrus.Infof("Extracting %d transports systems", len(transports_systems))

	for _, transport_system := range transports_systems {
		lines_transport_system, err := transport_system.Extract()
		if err != nil {
			return stations, []error{err}
		}

		lines = append(lines, lines_transport_system...)
	}

	logrus.Infof("Extracting %d lines", len(lines))

	for _, line := range lines {
		stations_line, err := line.Extract()
		if err != nil {
			return stations, []error{err}
		}

		stations = append(stations, stations_line...)
	}

	logrus.Infof("Extracting %d stations", len(stations))

	for i := range stations {
		logrus.Debugf("Extracting %d", i)
		err := stations[i].Extract()
		logrus.Debugf("Extracted %d", i)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return stations, errs
}
