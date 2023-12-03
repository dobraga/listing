package utils

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

// GetKeys returns an array of keys from the input map.
//
// The input parameter is a map with string keys and values of type T.
// The return type is an array of strings.
func GetKeys[T interface{} | string](input map[string]T) []string {
	var keys []string
	for k := range input {
		keys = append(keys, k)
	}
	return keys
}

// Contains checks if a given string is present in a slice of strings.
//
// Parameters:
// - s: the slice of strings to search through.
// - e: the string to look for in the slice.
//
// Returns:
// - bool: true if the string is found, false otherwise.
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Min returns the minimum value between two inputs.
//
// The function takes two parameters, x and y, which can be of type int64, int, float64, or float32.
// It returns the smaller of the two input values.
func Min[T int64 | int | float64 | float32](x, y T) T {
	if x > y {
		return y
	}
	return x
}

// GetFirst returns the first element of a given slice.
//
// It takes a slice of values, a URL, and a variable name as inputs.
// It returns the first element of the slice.
func GetFirst[T int64 | int | float64 | float32 | string](listaValores []T, url string, variable string) T {
	var value T
	switch qtdElements := len(listaValores); {
	case qtdElements > 1:
		log.Warn(fmt.Sprintf(`Property "%s" with %v %s`, url, listaValores, variable))
		value = listaValores[0]
	case qtdElements == 1:
		value = listaValores[0]
	}

	return value
}

// LoadURL loads the content of a webpage from the given URL and parses it as an HTML document.
//
// It takes a string parameter `url` which represents the URL of the webpage to be loaded.
// The function returns a pointer to an `html.Node` and an error. The `html.Node` represents the root
// of the parsed HTML document, and the error represents any error that occurred during the loading
// or parsing process.
func LoadURL(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	return html.Parse(r)
}

// RemoveDuplicates returns a new slice that contains only the unique elements from the input slice of strings.
//
// The function takes a single parameter:
//   - strList: a slice of strings
//
// It returns a new slice of strings that contains only the unique elements from the input slice.
func RemoveDuplicates(strList []string) []string {
	list := []string{}
	for _, item := range strList {
		if !Contains(list, item) {
			list = append(list, item)
		}
	}
	return list
}
