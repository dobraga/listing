package main

import (
	"fetch/pkg/controller/routes"
	"fetch/pkg/database"

	// "fetch/pkg/domain/station"
	"fetch/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	utils.LoadSettings()
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate()
	if err != nil {
		panic(err)
	}

	// err := station.CheckExistsExtractCreate(db)
	// if err != nil {
	// 	panic(err)
	// }

	r := gin.Default()
	routes.InitRoutes(&r.RouterGroup)

	port := viper.GetString("BACKEND_PORT")
	logrus.Infof("Started server on port %s", port)
	r.Run(":" + port)
}
