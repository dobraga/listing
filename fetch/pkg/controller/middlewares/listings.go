package middlewares

import (
	"fetch/pkg/database"
	"fetch/pkg/domain/property"
	"fetch/pkg/models"
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Listings(c *gin.Context) {
	dbAny, _ := c.Get("db")
	db, _ := dbAny.(database.Database)

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

	properties, err := property.Fetch(config, db)
	if err != nil {
		log.Error(err)
		c.JSON(500, err)
		return
	}

	if config.ReturnProperties {
		c.JSON(200, properties)
	} else {
		c.JSON(200, "Saved successfully")
	}
}
