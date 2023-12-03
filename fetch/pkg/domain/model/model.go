package model

import (
	"bytes"
	"encoding/json"
	"fetch/pkg/database"
	"fetch/pkg/models"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Predict(db database.Database, properties *[]models.Property) error {
	data := map[string]*[]models.Property{"data": properties}

	body, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}

	url := fmt.Sprintf("http://%s:%s/predict",
		viper.GetString("MODEL_HOST"), viper.GetString("MODEL_PORT"))

	log.Debugf("Getting '%s' in  '%s'", body, url)

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
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

	for i := range *properties {
		logrus.Debugf("Update predicted %f for %s", response[i], (*properties)[i].Url)
		(*properties)[i].PredictPrice = response[i]
	}

	return nil
}
