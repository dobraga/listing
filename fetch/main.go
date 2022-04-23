package main

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	LoadSettings()
	DB = Connect()

	SaveStations(DB)

	mapSites := viper.GetStringMap("sites")
	todosSites := GetKeys(mapSites)

	r := gin.Default()

	r.GET("/locations/:location", func(c *gin.Context) {
		location := c.Param("location")
		locations, err := ListLocations(location, "vivareal")

		if err != nil {
			c.JSON(400, err)
			return
		} else {
			c.JSON(200, locations)
			return
		}

	})

	r.GET("listings/:business_type/:listing_type/:city/:locationId/:neighborhood/:state/:stateAcronym/:zone", func(c *gin.Context) {

		var errs []error
		wg := new(sync.WaitGroup)
		channelErr := make(chan []error)

		for _, origin := range todosSites {
			wg.Add(1)

			go func(o string, c *gin.Context, d *gorm.DB, wg *sync.WaitGroup) {
				defer wg.Done()
				location := Location{
					Local: Local{
						City:         c.Param("city"),
						Zone:         c.Param("zone"),
						State:        c.Param("state"),
						LocationId:   c.Param("locationId"),
						Neighborhood: c.Param("neighborhood"),
						StateAcronym: c.Param("stateAcronym"),
					},
					BusinessType: c.Param("business_type"),
					ListingType:  c.Param("listing_type"),
					Origin:       o,
				}

				_, err := FetchListings(d, location)
				if err != nil {
					channelErr <- err
				}
			}(origin, c, DB, wg)
		}

		wg.Wait()
		close(channelErr)

		for err := range channelErr {
			errs = append(errs, err...)
		}

		if errs != nil {
			c.JSON(400, errs)
			return
		} else {
			c.JSON(200, "Saved successfully")
			return
		}
	})

	port := viper.GetString("PORT")
	r.Run(":" + port)
}
