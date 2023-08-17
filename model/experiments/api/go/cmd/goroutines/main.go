package main

import (
	"api_go_test/model"
	"flag"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	Input [][6]float32
}

type Pred struct {
	pos  int
	pred float32
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
		if err := c.BindJSON(&requestBody); err != nil {
			panic(err)
		}

		predictions := make(map[int]float32)
		predictionsChan := make(chan Pred)
		defer close(predictionsChan)

		for pos, x := range requestBody.Input {
			go func(pos int, x [6]float32) {
				predictionsChan <- Pred{pos, model.Predict(x)}
			}(pos, x)
		}

		for i := 0; i < len(requestBody.Input); i++ {
			p := <-predictionsChan
			predictions[p.pos] = p.pred
		}

		c.JSON(http.StatusOK, gin.H{
			"predictions": predictions,
		})
	})
	r.Run()
}
