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

type PageStation struct {
	Uf       string
	Title    string
	Linha    string
	Url      string
	URLLinha string
}

func SaveStations(db *gorm.DB) {
	pages := listStations()
	savePageStations(db, pages)
}

func savePageStations(db *gorm.DB, pages []PageStation) {

	wg := new(sync.WaitGroup)

	for _, page := range pages {
		wg.Add(1)

		go func(p PageStation, w *sync.WaitGroup) {
			defer w.Done()

			station := fetchStation(p)
			db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(station, 50)
		}(page, wg)
	}

	wg.Wait()

}

func listStations() []PageStation {
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

	return output
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

func fetchStation(page PageStation) Station {

	var lat, lon float64

	doc, err := htmlquery.LoadURL(page.Url)
	if err != nil {
		log.Error(err)
	}

	name := htmlquery.InnerText(htmlquery.FindOne(doc, "//*[@id='firstHeading']/text()"))

	urlLatlng := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//a[contains(@href, 'tools.wmflabs.org')]"), "href")

	if len(urlLatlng) > 0 {
		doc, err = htmlquery.LoadURL(urlLatlng)
		if err != nil {
			log.Error(err)
		}

		latlng := htmlquery.Find(doc, "//*[contains(@class, 'geo')]/span/text()")

		lat, err = strconv.ParseFloat(htmlquery.InnerText(latlng[0]), 64)
		if err != nil {
			log.Error(err)
		}
		lon, err = strconv.ParseFloat(htmlquery.InnerText(latlng[1]), 64)
		if err != nil {
			log.Error(err)
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

	return station
}
