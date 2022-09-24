package main

import (
	"fetch/database"
	"fetch/property"
	"fetch/station"
	"fetch/utils"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func main() {
	utils.LoadSettings()
	DB := database.Connect()
	database.AutoMigrate(DB)

	err := station.CheckExistsExtractCreate(DB)
	if err != nil {
		panic(err)
	}

	mapSites := viper.GetStringMap("sites")
	todosSites := utils.GetKeys(mapSites)

	r := gin.Default()

	r.GET("/locations/:location", func(c *gin.Context) {
		location := c.Param("location")
		locations, err := property.ListLocations(location, "vivareal")

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
		channelErr := make(chan []error)

		location := property.SearchConfig{
			Local: property.Local{
				City:         c.Param("city"),
				Zone:         c.Param("zone"),
				State:        c.Param("state"),
				LocationId:   c.Param("locationId"),
				Neighborhood: c.Param("neighborhood"),
				StateAcronym: c.Param("stateAcronym"),
			},
			BusinessType: c.Param("business_type"),
			ListingType:  c.Param("listing_type"),
			Origin:       "",
		}

		errs := location.Validation()
		if len(errs) > 0 {
			c.JSON(400, fmt.Sprintf("%v", errs))
			return
		}

		for _, origin := range todosSites {
			wg.Add(1)

			go func(o string, c *gin.Context, d *gorm.DB) {
				defer wg.Done()

				l := location
				l.Origin = o

				err := property.FetchProperties(DB, l, viper.GetInt("max_page"))
				if err != nil {
					channelErr <- err
				}
			}(origin, c, DB)
		}

		wg.Wait()
		close(channelErr)

		for err := range channelErr {
			errs = append(errs, err...)
		}

		if errs != nil {
			c.JSON(400, fmt.Sprintf("%v", errs))
			return
		} else {
			c.JSON(200, "Saved successfully")
			return
		}
	})

	port := viper.GetString("PORT")
	r.Run(":" + port)
}
