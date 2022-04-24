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
	var pageListings Property

	errs := location.Validation()
	if errs != nil {
		return "", errs
	}

	size := 24
	origin := location.Origin

	maxPage := viper.GetInt64("max_page")

	query := createQuery(location, size)

	qtdListings, err := qtdListings(origin, query)
	if err != nil {
		return "", []error{err}
	}
	total_pages := int64(qtdListings / size)
	if maxPage <= 0 {
		maxPage = total_pages
	} else {
		maxPage = Min(maxPage, total_pages)
	}

	log.Infof("Getting %d/%d pages with %d listings from '%s'", maxPage, total_pages, qtdListings, origin)

	wg := new(sync.WaitGroup)
	channelErr := make(chan error)

	for page := 1; page <= int(maxPage); page++ {
		wg.Add(1)
		log.Debugf("Getting page %d from '%s'", page, origin)
		query["from"] = page * query["size"].(int)

		bytesData := MakeRequest(false, origin, query)

		go func(p int, l Location, b []byte, d *gorm.DB, w *sync.WaitGroup, c chan error) {
			defer w.Done()

			listings, err := pageListings.Unmarshal(b, l)
			if err != nil {
				log.Errorf("Parsed error page %d from '%s': %v", p, l.Origin, err)
				c <- err
				return
			}

			result := d.Clauses(clause.OnConflict{DoNothing: true}).Create(listings)
			if result.Error != nil {
				log.Error(result.Error)
				c <- err
				return
			}
			log.Debugf("Saved successfully %d from '%s'", p, l.Origin)

		}(page, location, bytesData, DB, wg, channelErr)

		if page < int(maxPage) {
			time.Sleep(300 * time.Millisecond)
		}
	}

	wg.Wait()
	close(channelErr)

	for err = range channelErr {
		log.Info(err)
		errs = append(errs, err)
	}

	log.Infof("Saved %d pages from '%s'", maxPage, origin)

	return fmt.Sprintf("Saved %d pages from '%s'", maxPage, origin), errs
}

func qtdListings(origin string, query map[string]interface{}) (int, error) {

	bytesData := MakeRequest(false, origin, query)

	data := map[string]interface{}{}
	err := json.Unmarshal(bytesData, &data)
	if err != nil {
		err := fmt.Errorf(fmt.Sprintf(
			"erro ao buscar a quantidade de propriedades da página '%s' '%v' '%v': %v",
			origin, query, bytesData, err,
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

	qtdListings := data["totalCount"].(float64)

	return int(qtdListings), nil

}

func createQuery(location Location, size int) map[string]interface{} {
	return map[string]interface{}{
		"includeFields":       "search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount),expansion(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount)),nearby(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount)),page,fullUriFragments,developments(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount)),superPremium(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount)),owners(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount))",
		"addressNeighborhood": location.Local.Neighborhood,
		"addressLocationId":   location.Local.LocationId,
		"addressState":        location.Local.State,
		"addressCity":         location.Local.City,
		"addressZone":         location.Local.Zone,
		"listingType":         location.ListingType,
		"business":            location.BusinessType,
		"usageTypes":          "RESIDENTIAL",
		"categoryPage":        "RESULT",
		"size":                size,
		"from":                24,
	}
}
