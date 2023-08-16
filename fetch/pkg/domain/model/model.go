package model

import (
	"bytes"
	"encoding/json"
	"fetch/pkg/models"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Predict(p models.Property) (models.Property, error) {
	one := []models.Property{p}
	data := map[string][]models.Property{"data": one}

	body, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return p, err
	}
	fmt.Println(string(body))

	url := fmt.Sprintf(
		"http://%s:%s/%s", viper.GetString("MODEL_HOST"),
		viper.GetString("MODEL_PORT"), strings.ToLower(p.BusinessType))
	fmt.Println(url)

	log.Debugf("Getting '%s' in  '%s'", body, url)

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return p, err
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	var response []float64
	resBody, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(resBody, &response)

	p.PredictTotalPrice = response[0]

	return p, nil
}
