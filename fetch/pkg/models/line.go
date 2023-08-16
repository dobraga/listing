package models

import (
	"github.com/antchfx/htmlquery"
)

type Line struct {
	Name            string
	URL             string
	TransportSystem TransportSystem
}

func (l *Line) Extract() ([]Station, error) {
	var stations []Station

	doc, err := htmlquery.LoadURL(l.URL)
	if err != nil {
		return stations, err
	}

	node_estations := htmlquery.Find(doc, "//tr/td/a[contains(@href, 'wiki/Esta')]")

	for _, node := range node_estations {
		url := "https://pt.wikipedia.org" + htmlquery.SelectAttr(node, "href")

		stations = append(stations, Station{
			URL: url, LineName: l.Name, TransportSystemUf: l.TransportSystem.Uf,
			TransportSystemURL: l.TransportSystem.URL})
	}

	return stations, nil
}
