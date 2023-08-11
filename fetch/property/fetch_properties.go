package property

import (
	"fetch/database"
	"fetch/models"
	"fetch/utils"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

func StoreFetchProperties(config models.SearchConfig, maxPage int) error {
	properties, err := FetchProperties(config, maxPage)
	if err != nil {
		return err
	}

	return database.StoreProperty(config, properties)
}

func FetchProperties(config models.SearchConfig, maxPage int) ([]models.Property, error) {
	log.Infof("Searching listings from %+v", config)
	var properties []models.Property
	size := 24
	query := createQuery(config, size)

	qtdListings, err := qtdListings(config, query)
	if err != nil {
		return properties, err
	}
	total_pages := int(qtdListings / size)
	if maxPage < 0 {
		maxPage = total_pages
	} else {
		maxPage = utils.Min(maxPage, total_pages)
	}
	log.Infof("Getting %d/%d pages with %d total listings from '%s'", maxPage, total_pages, qtdListings, config.Origin)

	for page := 0; page <= maxPage; page++ {
		log.Debugf("Getting page %d from '%s'", page, config.Origin)
		query["from"] = page * query["size"].(int)

		data, err := MakeRequest(false, config.Origin, query)
		if err != nil {
			log.Error(err)
			continue
		}

		property, err := UnmarshalProperty(data, config)
		if err != nil {
			log.Error(err)
			continue
		}

		log.Infof("add %d properties from page %d '%s'", len(property), page, config.Origin)
		properties = append(properties, property...)

		if page < maxPage {
			log.Debugf("Sleeping after extract %d page", page)
			time.Sleep(2 * time.Second)
		}
	}

	return properties, err
}

func qtdListings(config models.SearchConfig, query map[string]interface{}) (int, error) {
	data, err := MakeRequest(false, config.Origin, query)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	if utils.Contains(utils.GetKeys(data), "nearby") {
		data = data["nearby"].(map[string]interface{})
	}

	if !utils.Contains(utils.GetKeys(data), "search") {
		err := fmt.Errorf("not found search listings '%v' from '%s' '%v'", data, config.Origin, query)
		log.Error(err)
		return 0, err
	}

	data = data["search"].(map[string]interface{})
	qtd := data["totalCount"]
	return int(qtd.(float64)), nil
}

func createQuery(config models.SearchConfig, size int) map[string]interface{} {
	data := map[string]interface{}{
		"addressNeighborhood": config.Local.Neighborhood,
		"addressLocationId":   config.Local.LocationId,
		"addressState":        config.Local.State,
		"addressCity":         config.Local.City,
		"addressZone":         config.Local.Zone,
		"listingType":         config.ListingType,
		"business":            config.BusinessType,
		"usageTypes":          "RESIDENTIAL",
		"categoryPage":        "RESULT",
		"size":                size,
		"from":                24,
	}

	if config.Local.AddressStreet != "" {
		data["addressStreet"] = config.Local.AddressStreet
		data["addressPointLat"] = config.Local.AddressPointLat
		data["addressPointLon"] = config.Local.AddressPointLon
	}

	log.Debugf("query: %v", data)
	data["includeFields"] = "search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount),expansion(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount)),nearby(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount)),page,fullUriFragments,developments(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount)),superPremium(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount)),owners(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount))"

	return data
}
