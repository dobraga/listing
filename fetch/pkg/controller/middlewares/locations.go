package middlewares

import (
	"fetch/pkg/domain/property"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	locations, status_code, err := property.ListLocations(location, type_location, "vivareal")

	if err != nil {
		logrus.Error(err)
		c.JSON(status_code, err)
		return
	} else {
		c.JSON(200, locations)
		return
	}
}
