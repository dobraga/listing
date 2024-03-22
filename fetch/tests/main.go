package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Given URL and headers
	url := "https://glue-api.vivareal.com/v3/locations?q=Rua%20Bar%C3%A3o&portal=VIVAREAL&size=6&fields=neighborhood,city,street&includeFields=address.city,address.zone,address.state,address.neighborhood,address.stateAcronym,address.street,address.locationId,address.point&"
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 ",
	}

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set request headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response body
	fmt.Println("Response body:", string(body))
	fmt.Print(req.Header)
}
