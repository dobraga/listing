package main

import (
	"fetch/database"
	"fetch/models"
	"fetch/property"
	"fetch/station"
	"fetch/utils"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	utils.LoadSettings()
	database.AutoMigrate()

	err := station.CheckExistsExtractCreate()
	if err != nil {
		panic(err)
	}

	mapSites := viper.GetStringMap("sites")
	todosSites := utils.GetKeys(mapSites)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/locations", func(c *gin.Context) {
		type_location := c.DefaultQuery("type", "neighborhood")
		location := c.Query("q")
		if location == "" {
			location = c.Query("location")
		}

		if location == "" {
			c.JSON(400, "need a non empty string into 'location' or 'q'")
			return
		}
		locations, err := property.ListLocations(location, type_location, "vivareal")

		if err != nil {
			c.JSON(500, err)
			return
		} else {
			c.JSON(200, locations)
			return
		}
	})

	r.GET("listings", func(c *gin.Context) {
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
			err_str := fmt.Sprintf("Errors: %+v", errs)
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

				err := property.StoreFetchProperties(l, viper.GetInt("max_page"))
				if err != nil {
					channelErr <- err
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
			c.JSON(200, "Saved successfully")
			return
		}
	})

	r.GET("get_listings", func(c *gin.Context) {
		var all_properties []models.Property

		location := models.SearchConfig{}
		location.ExtractFromContext(c)
		errs := location.Validation()
		if len(errs) > 0 {
			err_str := fmt.Sprintf("Errors: %+v", errs)
			logrus.Errorf(err_str)
			c.JSON(500, err_str)
			return
		}

		for _, origin := range todosSites {
			l := location
			l.Origin = origin

			properties, err := property.FetchProperties(l, viper.GetInt("max_page"))
			if err != nil {
				errs = append(errs, err)
			}
			all_properties = append(all_properties, properties...)
		}

		if len(errs) > 0 {
			err_str := fmt.Sprintf("Errors: %+v", errs)
			logrus.Errorf(err_str)
			c.JSON(500, err_str)
			return
		} else {
			c.JSON(200, all_properties)
			return
		}
	})

	port := viper.GetString("BACKEND_PORT")
	r.Run(":" + port)
}
