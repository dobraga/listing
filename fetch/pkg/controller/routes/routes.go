package routes

import (
	"fetch/pkg/controller/middlewares"
	"fetch/pkg/database"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.RouterGroup, db database.Database) {
	r.GET("/health", middlewares.Health)
	r.GET("/locations", middlewares.ListLocation)
	r.GET("/listings", middlewares.Conf(db), middlewares.Listings)
}
