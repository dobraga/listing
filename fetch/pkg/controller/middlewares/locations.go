package middlewares

import (
	"fetch/pkg/domain/property"

	"github.com/gin-gonic/gin"
)

func ListLocation(c *gin.Context) {
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
}
