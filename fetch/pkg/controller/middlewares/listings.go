package middlewares

import (
	"fetch/pkg/database"
	"fetch/pkg/domain/property"
	"fetch/pkg/models"
	"fetch/pkg/utils"
	"fmt"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func StoreListings(c *gin.Context) {
	var all_properties []models.Property
	var statusCode int
	db, err := database.Connect()
	if err != nil {
		logrus.Error(err)
		c.JSON(500, err)
		return
	}
	defer db.Close()

	returnListings, err := strconv.ParseBool(c.DefaultQuery("return_listings", "false"))
	if err != nil {
		logrus.Error(err)
		c.JSON(500, err)
		return
	}
	location := models.SearchConfig{}
	location.ExtractFromContext(c)
	errs := location.Validation()
	if len(errs) > 0 {
		err_str := fmt.Sprintf("Errors: %+v", errs)
		logrus.Errorf(err_str)
		c.JSON(500, err_str)
		return
	}

	if location.StoreProperty {
		err = db.ResetActive(location)
		if err != nil {
			err_str := fmt.Sprintf("Error: %+v", err)
			logrus.Errorf(err_str)
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

		go func(o string, c *gin.Context) {
			defer wg.Done()

			l := location
			l.Origin = o

			properties, statusCode, err := property.FetchProperties(l)
			if err != nil {
				channelErr <- err
				channelStatusCode <- statusCode
				return
			}

			if l.StoreProperty {
				err = db.StoreProperty(l, properties)
				if err != nil {
					channelErr <- err
					return
				}
			}
			if returnListings {
				all_properties = append(all_properties, properties...)
			}
		}(origin, c)
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
		logrus.Errorf(err_str)
		c.JSON(statusCode, err_str)
		return
	} else {
		if returnListings {
			c.JSON(200, all_properties)
		} else {
			c.JSON(200, "Saved successfully")
		}
		return
	}
}
