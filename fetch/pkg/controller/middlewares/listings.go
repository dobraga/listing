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
	returnListings, err := strconv.ParseBool(c.DefaultQuery("return_listings", "false"))
	if err != nil {
		err_str := fmt.Sprintf("Error: %+v", err)
		logrus.Errorf(err_str)
		c.JSON(500, err_str)
		return
	}

	var all_properties []models.Property

	mapSites := viper.GetStringMap("sites")
	todosSites := utils.GetKeys(mapSites)

	wg := new(sync.WaitGroup)
	channelErr := make(chan error)

	location := models.SearchConfig{}
	location.ExtractFromContext(c)
	errs := location.Validation()
	if len(errs) > 0 {
		err_str := fmt.Sprintf("Errors: %+v", errs)
		logrus.Errorf(err_str)
		c.JSON(500, err_str)
		return
	}

	err = database.ResetActive(location)
	if err != nil {
		err_str := fmt.Sprintf("Error: %+v", err)
		logrus.Errorf(err_str)
		c.JSON(500, err_str)
		return
	}

	for _, origin := range todosSites {
		wg.Add(1)

		go func(o string, c *gin.Context) {
			defer wg.Done()

			l := location
			l.Origin = o

			properties, err := property.FetchProperties(l, viper.GetInt("max_page"))
			if err != nil {
				channelErr <- err
				return
			}

			err = database.StoreProperty(location, properties)
			if err != nil {
				channelErr <- err
				return
			}
			if returnListings {
				all_properties = append(all_properties, properties...)
			}
		}(origin, c)
	}

	wg.Wait()
	close(channelErr)

	for err := range channelErr {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		err_str := fmt.Sprintf("Errors: %+v", errs)
		logrus.Errorf(err_str)
		c.JSON(500, err_str)
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
