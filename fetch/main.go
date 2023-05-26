package main

import (
	"fetch/database"
	"fetch/models"
	"fetch/property"
	"fetch/station"
	"fetch/utils"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
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

	r.GET("/locations/:location", func(c *gin.Context) {
		location := c.Param("location")
		locations, err := property.ListLocations(location, "neighborhood", "vivareal")

		if err != nil {
			c.JSON(400, err)
			return
		} else {
			c.JSON(200, locations)
			return
		}
	})

	r.GET("listings/:business_type/:listing_type/:city/:locationId/:neighborhood/:state/:stateAcronym/:zone", func(c *gin.Context) {
		wg := new(sync.WaitGroup)
		channelErr := make(chan error)

		location := models.SearchConfig{}
		location.ExtractFromContext(c)
		errs := location.Validation()
		if len(errs) > 0 {
			c.JSON(400, errs)
			return
		}

		err = database.ResetActive(location)
		if err != nil {
			c.JSON(400, err)
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

		if errs != nil {
			c.JSON(400, errs)
			return
		} else {
			c.JSON(200, "Saved successfully")
			return
		}
	})

	port := viper.GetString("BACKEND_PORT")
	r.Run(":" + port)
}
