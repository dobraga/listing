package middlewares

import (
	"fetch/pkg/database"

	"github.com/gin-gonic/gin"
)

func Conf(db database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
