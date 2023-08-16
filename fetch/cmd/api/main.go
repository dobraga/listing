package main

import (
	"fetch/pkg/controller/routes"
	"fetch/pkg/database"
	"fetch/pkg/domain/station"
	"fetch/pkg/utils"

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

	r := gin.Default()
	routes.InitRoutes(&r.RouterGroup)

	port := viper.GetString("BACKEND_PORT")
	r.Run(":" + port)
}
