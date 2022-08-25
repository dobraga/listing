package main

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"

	"github.com/antchfx/htmlquery"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var re = regexp.MustCompile("[A-zÀ-ú ]+")

func SaveStations(db *gorm.DB) {
	pages, err := ListStations()
	var stations []Station

	if err != nil {
		wg := new(sync.WaitGroup)

		for _, page := range pages {
			wg.Add(1)

			go func(p PageStation, w *sync.WaitGroup) {
				defer w.Done()

				station, err := FetchStationRetry(p)
				if err != nil {
					stations = append(stations, station)
				}
			}(page, wg)
		}

		wg.Wait()

		db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(stations, 50)
	}
}

func ListStations() ([]PageStation, error) {
	var output []PageStation

	urls := readUrlsSettings()

	log.Debugf("Fetching: %v", urls)

	for _, ufUrls := range urls {
		uf := ufUrls[0]
		url := ufUrls[1]
		log.Infof("Getting %s '%s'", uf, url)
		linhas := map[string]string{}
		urls := []string{}

		doc, err := htmlquery.LoadURL(url)
		if err != nil {
			log.Error(err)
			return []PageStation{}, err
		}

		title := htmlquery.InnerText(htmlquery.FindOne(doc, "//*[@id='firstHeading']/text()"))
		list := htmlquery.Find(doc, "(//tr/th[contains(text(), 'erminais')]/../..)[1]//td[1]//a")

		for _, l := range list {
			linha := fmt.Sprintf("https://pt.wikipedia.org%s", htmlquery.SelectAttr(l, "href"))
			name := htmlquery.InnerText(l)
			name = re.FindStringSubmatch(name)[0]

			linhas[name] = linha
		}

		for linha, URLLinha := range linhas {
			log.Debugf("Getting %s '%s' '%s'", uf, linha, url)
			doc, err := htmlquery.LoadURL(URLLinha)
			if err != nil {
				log.Error(err)
				return []PageStation{}, err
			}

			list := htmlquery.Find(doc, "//tr/td/a[contains(@href, 'wiki/Esta')]")

			for _, l := range list {
				url := fmt.Sprintf("https://pt.wikipedia.org%s", htmlquery.SelectAttr(l, "href"))
				urls = append(urls, url)
			}

			urls = RemoveDuplicates(urls)

			for _, l := range list {
				url := fmt.Sprintf("https://pt.wikipedia.org%s", htmlquery.SelectAttr(l, "href"))
				output = append(output, PageStation{uf, title, linha, url, URLLinha})
			}
		}
	}

	return output, nil
}

func FetchStationRetry(page PageStation) (Station, error) {
	var err error
	var station Station

	for i := 0; i < 5; i++ {
		station, err = fetchStation(page)

		if err == nil {
			return station, nil
		}

		log.Infof("Retrying %d", i)
	}

	return Station{}, err
}

func fetchStation(page PageStation) (Station, error) {

	var lat, lon float64

	doc, err := htmlquery.LoadURL(page.Url)
	if err != nil {
		log.Error(err)
		return Station{}, err
	}

	name := htmlquery.InnerText(htmlquery.FindOne(doc, "//*[@id='firstHeading']/text()"))

	urlLatlng := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//a[contains(@href, 'tools.wmflabs.org')]"), "href")

	if len(urlLatlng) > 0 {
		doc, err = htmlquery.LoadURL(urlLatlng)
		if err != nil {
			log.Error(err)
			return Station{}, err
		}

		latlng := htmlquery.Find(doc, "//*[@class = 'geo h-geo']/span/text()")

		if len(latlng) != 2 {
			err = fmt.Errorf("latlong wrong scrap into '%s', find %d elements",
				urlLatlng, len(latlng))
			log.Error(err)
			return Station{}, err
		}

		str_lat := htmlquery.InnerText(latlng[0])
		lat, err = strconv.ParseFloat(str_lat, 64)
		if err != nil {
			err = fmt.Errorf(
				"cannot convert %s to float, into page '%s'",
				str_lat, urlLatlng)
			log.Error(err)
			return Station{}, err
		}
		str_lon := htmlquery.InnerText(latlng[1])
		lon, err = strconv.ParseFloat(str_lon, 64)
		if err != nil {
			err = fmt.Errorf(
				"cannot convert %s to float, into page '%s'",
				str_lon, urlLatlng)
			log.Error(err)
			return Station{}, err
		}
	}

	station := Station{
		Name:     page.Title,
		Station:  name,
		Uf:       page.Uf,
		Linha:    page.Linha,
		Lat:      lat,
		Lon:      lon,
		Url:      page.Url,
		URLLinha: page.URLLinha,
	}

	return station, nil
}

func readUrlsSettings() [][2]string {
	urlMetroTrem := viper.GetStringMap("metro_trem")
	allUrls := [][2]string{}

	for uf, urls := range urlMetroTrem {
		mapUrls := urls.(map[string]interface{})

		for _, url := range mapUrls["urls"].([]interface{}) {
			allUrls = append(allUrls, [2]string{uf, url.(string)})
		}
	}

	return allUrls
}
