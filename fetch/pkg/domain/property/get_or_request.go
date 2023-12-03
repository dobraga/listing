package property

import (
	"fetch/pkg/database"
	"fetch/pkg/domain/model"
	"fetch/pkg/models"
	"fetch/pkg/utils"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Fetch(config models.SearchConfig, db database.Database) ([]models.Property, error) {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	var properties []models.Property

	lastUpdate, err := db.GetLastUpdate(config)
	if err == nil {
		lastUpdate = lastUpdate.In(loc)
		now := time.Now().In(loc)
		lastUpdateHrs := now.Sub(lastUpdate).Hours()

		log.Infof("Last update %f hours ago: %s",
			lastUpdateHrs, lastUpdate.Format(time.RFC3339))

		if lastUpdateHrs <= 3 {
			return db.Where(config)
		}
	} else {
		log.Error(err)
	}

	err = db.ResetActive(config)
	if err != nil {
		return properties, err
	}
	properties, err = getNewProperties(db, config)
	if err != nil {
		db.ResetActive(config)
		return properties, err
	}

	err = model.Predict(db, &properties)
	if err != nil {
		return properties, err
	}

	return properties, db.StoreProperty(config, properties)
}

func getNewProperties(db database.Database, config models.SearchConfig) ([]models.Property, error) {
	var allProperties []models.Property
	var errs []error

	mapSites := viper.GetStringMap("sites")
	todosSites := utils.GetKeys(mapSites)

	wg := new(sync.WaitGroup)
	channelErr := make(chan error)
	channelStatusCode := make(chan int)

	for _, origin := range todosSites {
		wg.Add(1)

		go func(o string, c models.SearchConfig) {
			defer wg.Done()

			c.Origin = o

			properties, statusCode, err := FetchProperties(c)
			if err != nil {
				channelErr <- err
				channelStatusCode <- statusCode
				return
			}
			allProperties = append(allProperties, properties...)
		}(origin, config)
	}

	wg.Wait()
	close(channelErr)

	for err := range channelErr {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return allProperties, fmt.Errorf("errors: %+v", errs)
	}

	return allProperties, nil
}
