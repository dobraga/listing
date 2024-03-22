package property

import (
	"encoding/json"
	"fmt"
	"io"
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

	if origin == "" {
		return nil, fmt.Errorf("origin cannot be empty")
	}

	siteInfo := viper.Get("sites").(map[string]interface{})[origin].(map[string]interface{})
	if location {
		url = fmt.Sprintf("https://%s/v3/locations", siteInfo["api"])
	} else {
		url = fmt.Sprintf("https://%s/v2/listings", siteInfo["api"])
	}

	req, err := CreateRequestWithQueryHeaders(url, query)
	if err != nil {
		return nil, err
	}

	// Request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição da página '%v': %v", req.URL, err)
	}
	defer resp.Body.Close()

	// Response to interface
	bytesData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body for URL '%v': %v", req.URL, err)
	}

	// verify status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro na requisição da página '%v': status code %v, body: %s", req.URL, resp.StatusCode, string(bytesData))
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

func CreateRequestWithQueryHeaders(url string, query map[string]interface{}) (*http.Request, error) {
	queryString := "?"
	for key, element := range query {
		queryString += fmt.Sprintf("%s=%v&", key, element)
	}
	url += queryString

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição da página '%s': %v", url, err)
	}

	// add headers
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0",
	}
	logrus.Infof("Using headers %v+", headers)
	for key, element := range headers {
		req.Header.Add(key, fmt.Sprintf("%v", element))
	}
	logrus.Infof("Requisição da pagina '%s", req.URL)

	return req, nil
}
