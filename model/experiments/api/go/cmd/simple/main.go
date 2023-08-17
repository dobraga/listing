package main

import (
	"flag"
	"net/http"

	"github.com/gin-gonic/gin"

	m "api_go_test/model"
)

type RequestBody struct {
	Input [][6]float32
}

func main() {
	prod := flag.Bool("prod", false, "Define production mode")
	flag.Parse()

	if *prod {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.POST("/predict", func(c *gin.Context) {
		var requestBody RequestBody
		preds := []float32{}

		if err := c.BindJSON(&requestBody); err != nil {
			panic(err)
		}

		for _, x := range requestBody.Input {
			pred := m.Predict(x)
			preds = append(preds, pred)
		}

		c.JSON(http.StatusOK, gin.H{
			"predict": preds,
		})
	})
	r.Run()
}
