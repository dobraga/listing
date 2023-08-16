package models

import (
	"github.com/antchfx/htmlquery"
)

type TransportSystem struct {
	URL  string
	Uf   string
	Name string
}

func (t *TransportSystem) Extract() ([]Line, error) {
	var lines []Line

	doc, err := htmlquery.LoadURL(t.URL)
	if err != nil {
		return lines, err
	}

	t.Name = htmlquery.InnerText(htmlquery.FindOne(doc, "//*[@id='firstHeading']"))
	terminals := htmlquery.Find(doc, "(//tr/th[contains(text(), 'erminais')]/../..)[1]//td[1]//a")

	for _, l := range terminals {
		line := Line{
			URL:             "https://pt.wikipedia.org" + htmlquery.SelectAttr(l, "href"),
			Name:            htmlquery.InnerText(l),
			TransportSystem: *t,
		}

		lines = append(lines, line)
	}

	return lines, nil
}
