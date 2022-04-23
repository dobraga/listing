package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

func GetKeys[T interface{} | string](input map[string]T) []string {
	var keys []string
	for k := range input {
		keys = append(keys, k)
	}
	return keys
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Min(x, y int64) int64 {
	if x > y {
		return y
	}
	return x
}

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

func RemoveDuplicates(strList []string) []string {
	list := []string{}
	for _, item := range strList {
		if !Contains(list, item) {
			list = append(list, item)
		}
	}
	return list
}
