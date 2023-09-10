package middlewares

import (
	"fetch/pkg/database"
	"fetch/pkg/domain/property"
	"fetch/pkg/models"
	"fetch/pkg/utils"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func StoreListings(c *gin.Context) {
	var all_properties []models.Property
	var statusCode int
	db, err := database.Connect()
	if err != nil {
		log.Error(err)
		c.JSON(500, err)
		return
	}
	defer db.Close()

	config := models.SearchConfig{}
	config.ExtractFromContext(c)
	errs := config.Validation()
	if len(errs) > 0 {
		err_str := fmt.Sprintf("Errors: %+v", errs)
		log.Errorf(err_str)
		c.JSON(500, err_str)
		return
	}
	log.Infof("Searching listings from %s", config.String())

	if config.StoreProperty {
		lastUpdate, err := db.GetLastUpdate(config)
		if err == nil {
			now := time.Now()
			log.Infof("Last update: %s", lastUpdate.Format(time.RFC3339))

			if lastUpdate.Day() == now.Day() && lastUpdate.Month() == now.Month() && lastUpdate.Year() == now.Year() {
				c.JSON(200, "already updated at"+lastUpdate.Format(time.RFC3339))
				return
			}
		} else {
			log.Error(err)
		}

		err = db.ResetActive(config)
		if err != nil {
			err_str := fmt.Sprintf("Error: %+v", err)
			log.Errorf(err_str)
			c.JSON(500, err_str)
			return
		}
	}

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

			properties, statusCode, err := property.FetchProperties(c)
			if err != nil {
				channelErr <- err
				channelStatusCode <- statusCode
				return
			}

			if c.StoreProperty {
				err = db.StoreProperty(c, properties)
				if err != nil {
					channelErr <- err
					return
				}
			}
			if c.ReturnListings {
				all_properties = append(all_properties, properties...)
			}
		}(origin, config)
	}

	wg.Wait()
	close(channelErr)
	close(channelStatusCode)

	for err := range channelErr {
		errs = append(errs, err)
	}
	for statusCode = range channelStatusCode {
		continue
	}

	if len(errs) > 0 {
		err_str := fmt.Sprintf("Errors: %+v", errs)
		log.Errorf(err_str)
		c.JSON(statusCode, err_str)
		return
	} else {
		if config.ReturnListings {
			c.JSON(200, all_properties)
		} else {
			c.JSON(200, "Saved successfully")
		}
		return
	}
}
