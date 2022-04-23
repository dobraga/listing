package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var client http.Client = http.Client{}

func MakeRequest(location bool, origin string, query map[string]interface{}) []byte {
	var url string

	site_info := viper.Get("sites").(map[string]interface{})[origin].(map[string]interface{})
	if location {
		url = fmt.Sprintf("https://%s/v3/locations", site_info["api"])
	} else {
		url = fmt.Sprintf("https://%s/v2/listings", site_info["api"])
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(fmt.Sprintf("Erro na requisição da página '%s': %v", url, err))
	}

	// Query String
	q := req.URL.Query()
	for key, element := range query {
		q.Add(key, fmt.Sprintf("%v", element))
	}

	req.URL.RawQuery = q.Encode()

	// Headers
	headers := makeHeaders(origin)
	for key, element := range headers {
		req.Header.Add(key, fmt.Sprintf("%v", element))
	}

	// Request
	resp, err := client.Do(req)
	if err != nil {
		log.Error(fmt.Sprintf("Erro na requisição da página '%s' %v: %v", url, query, err))
	}
	defer resp.Body.Close()

	// Response to interface
	bytes_data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(fmt.Sprintf("Erro no parse da página '%s' %v: %v", url, query, err))
	}

	return bytes_data
}

func makeHeaders(origin string) map[string]string {
	return map[string]string{
		"user-agent":       "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
		"accept-language":  "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7",
		"sec-fetch-site":   "same-site",
		"accept":           "application/json",
		"sec-fetch-dest":   "empty",
		"sec-ch-ua-mobile": "?0",
		"sec-fetch-mode":   "cors",
		"origin-ua-mobile": "?0",
		"referer":          fmt.Sprintf("https://www.%s.com.br", origin),
		"origin":           fmt.Sprintf("https://www.%s.com.br", origin),
		"x-domain":         fmt.Sprintf("www.%s.com.br", origin),
	}
}
