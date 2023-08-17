package routes

import (
	"fetch/pkg/controller/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.RouterGroup) {
	r.GET("/health", middlewares.Health)
	r.GET("/locations", middlewares.ListLocation)
	r.GET("/listings", middlewares.StoreListings)
}
