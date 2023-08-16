package property

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var client http.Client = http.Client{}

// MakeRequest makes a HTTP GET request to a specific location or listings, with optional query parameters.
//
// Parameters:
// - location: a boolean indicating whether to request locations or listings.
// - origin: a string representing the origin to make the request to.
// - query: a map[string]interface{} containing optional query parameters.
//
// Returns:
// - map[string]interface{}: the response data as a map.
// - error: an error if the request fails.
func MakeRequest(location bool, origin string, query map[string]interface{}) (map[string]interface{}, error) {
	var url string
	var err error

	siteInfo := viper.Get("sites").(map[string]interface{})[origin].(map[string]interface{})
	if location {
		url = fmt.Sprintf("https://%s/v3/locations", siteInfo["api"])
	} else {
		url = fmt.Sprintf("https://%s/v2/listings", siteInfo["api"])
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição da página '%s': %v", url, err)
	}

	// Query String
	q := req.URL.Query()
	for key, element := range query {
		q.Add(key, fmt.Sprintf("%v", element))
	}

	req.URL.RawQuery = q.Encode()
	logrus.Infof("Requisição da pagina '%s", req.URL)

	// Headers
	headers := makeHeaders(origin)
	for key, element := range headers {
		req.Header.Add(key, fmt.Sprintf("%v", element))
	}

	// Request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição da página '%v': %v", req.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro na requisição da página '%v': status code %v", req.URL, resp.StatusCode)
	}

	// Response to interface
	bytesData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro no parse da página '%v': %v", req.URL, err)
	}

	// Bytes to map
	data := map[string]interface{}{}
	err = json.Unmarshal(bytesData, &data)
	if err != nil {
		return nil, fmt.Errorf("erro no parse da página '%v': %v", req.URL, err)
	}

	erro_value, ok := data["err"]
	if ok {
		return nil, fmt.Errorf("erro na requisição da página '%v': %v", req.URL, erro_value)
	}

	return data, nil
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
