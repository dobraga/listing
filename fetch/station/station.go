package station

import (
	"fmt"
	"strconv"

	"github.com/antchfx/htmlquery"
	"gorm.io/gorm"
)

type Station struct {
	gorm.Model
	URL                string `gorm:"primaryKey"`
	Name               string
	Lat                float64
	Lon                float64
	LineName           string
	TransportSystemUf  string
	TransportSystemURL string
}

func (station *Station) Extract() error {
	var lat, lon float64

	doc, err := htmlquery.LoadURL(station.URL)
	if err != nil {
		return err
	}

	station.Name = htmlquery.InnerText(htmlquery.FindOne(doc, "//*[@id='firstHeading']"))

	urlLatlng := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//a[contains(@href, 'tools.wmflabs.org')]"), "href")

	if len(urlLatlng) > 0 {
		doc, err = htmlquery.LoadURL(urlLatlng)
		if err != nil {
			return err
		}

		latlng := htmlquery.Find(doc, "//*[@class = 'geo h-geo']/span/text()")

		if len(latlng) != 2 {
			err = fmt.Errorf("latlong wrong scrap into '%s', find %d elements %v",
				urlLatlng, len(latlng), doc)

			return err
		}

		str_lat := htmlquery.InnerText(latlng[0])
		lat, err = strconv.ParseFloat(str_lat, 64)
		if err != nil {
			err = fmt.Errorf(
				"cannot convert %s to float, into page '%s'",
				str_lat, urlLatlng)
			return err
		}
		str_lon := htmlquery.InnerText(latlng[1])
		lon, err = strconv.ParseFloat(str_lon, 64)
		if err != nil {
			err = fmt.Errorf(
				"cannot convert %s to float, into page '%s'",
				str_lon, urlLatlng)
			return err
		}
		station.Lat = lat
		station.Lon = lon

		return nil
	}

	return fmt.Errorf("not found url to extract lat long '%s': '%s'", station.URL, station.Name)
}
