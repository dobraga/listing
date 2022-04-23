package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func FetchListings(DB *gorm.DB, location Location) (string, []error) {
	var page_listings Property

	errs := location.FinalValidation()
	if errs != nil {
		return "", errs
	}

	size := 24
	origin := location.Origin

	max_page := viper.GetInt64("max_page")

	query := createQuery(location, size)

	qtd_listings, err := qtdListings(origin, query)
	if err != nil {
		return "", []error{err}
	}
	total_pages := int64(qtd_listings / size)
	if max_page <= 0 {
		max_page = total_pages
	} else {
		max_page = Min(max_page, total_pages)
	}

	log.Infof("Getting %d/%d pages with %d listings from '%s'", max_page, total_pages, qtd_listings, origin)

	wg := new(sync.WaitGroup)
	channel_err := make(chan error)

	for page := 1; page <= int(max_page); page++ {
		wg.Add(1)
		log.Debugf("Getting page %d from '%s'", page, origin)
		query["from"] = page * query["size"].(int)

		bytes_data := MakeRequest(false, origin, query)

		go func(p int, o, bt string, b []byte, d *gorm.DB, w *sync.WaitGroup, c chan error) {
			defer w.Done()

			l, err := page_listings.Unmarshal(b, o, bt)
			if err != nil {
				log.Errorf("Parsed error page %d from '%s': %v", p, o, err)
				c <- err
				return
			}

			result := d.Clauses(clause.OnConflict{DoNothing: true}).Create(l)
			if result.Error != nil {
				log.Error(result.Error)
				c <- err
				return
			}
			log.Debugf("Saved successfully %d from '%s'", p, o)

		}(page, origin, location.BusinessType, bytes_data, DB, wg, channel_err)

		if page < int(max_page) {
			time.Sleep(300 * time.Millisecond)
		}
	}

	wg.Wait()
	close(channel_err)

	for err = range channel_err {
		log.Info(err)
		errs = append(errs, err)
	}

	log.Infof("Saved %d pages from '%s'", max_page, origin)

	return fmt.Sprintf("Saved %d pages from '%s'", max_page, origin), errs
}

func qtdListings(origin string, query map[string]interface{}) (int, error) {

	bytes_data := MakeRequest(false, origin, query)

	data := map[string]interface{}{}
	err := json.Unmarshal(bytes_data, &data)
	if err != nil {
		err := fmt.Errorf(fmt.Sprintf(
			"erro ao buscar a quantidade de propriedades da página '%s' '%v' '%v': %v",
			origin, query, bytes_data, err,
		))
		log.Error(err)
		return 0, err
	}

	if !Contains(GetKeys(data), "search") {
		err := fmt.Errorf("not found search listings '%v' from '%s' '%v'", data, origin, query)
		log.Error(err)
		return 0, err
	}

	data = data["search"].(map[string]interface{})

	qtd_listings := data["totalCount"].(float64)

	return int(qtd_listings), nil

}

func createQuery(location Location, size int) map[string]interface{} {
	return map[string]interface{}{
		"includeFields":       "search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount),expansion(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount)),nearby(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount)),page,fullUriFragments,developments(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount)),superPremium(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount)),owners(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount))",
		"addressNeighborhood": location.Neighborhood,
		"addressLocationId":   location.LocationId,
		"addressState":        location.State,
		"addressCity":         location.City,
		"addressZone":         location.Zone,
		"listingType":         location.ListingType,
		"business":            location.BusinessType,
		"usageTypes":          "RESIDENTIAL",
		"categoryPage":        "RESULT",
		"size":                size,
		"from":                24,
	}
}
