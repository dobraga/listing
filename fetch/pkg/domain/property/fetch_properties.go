package property

import (
	"fetch/pkg/models"
	"fetch/pkg/utils"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

func FetchProperties(config models.SearchConfig) ([]models.Property, int, error) {
	log.Infof("Searching listings from %s", config.String())
	var properties []models.Property
	size := 24
	query := createQuery(config, size)

	qtd_listings, status_code, err := qtdListings(config, query)
	if err != nil {
		return properties, status_code, err
	}
	total_pages := int(qtd_listings / size)
	maxPage := config.MaxPages
	if maxPage < 0 {
		maxPage = total_pages
	} else {
		maxPage = utils.Min(maxPage, total_pages)
	}
	log.Infof("Getting %d/%d pages with %d total listings from '%s'", maxPage, total_pages, qtd_listings, config.Origin)

	for page := 0; page <= maxPage; page++ {
		log.Debugf("Getting page %d from '%s'", page, config.Origin)
		query["from"] = page * query["size"].(int)

		data, _, err := MakeRequest(false, config.Origin, query)
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

	return properties, status_code, err
}

func qtdListings(config models.SearchConfig, query map[string]interface{}) (int, int, error) {
	data, status_code, err := MakeRequest(false, config.Origin, query)
	if err != nil {
		log.Error(err)
		return 0, status_code, err
	}

	if utils.Contains(utils.GetKeys(data), "nearby") {
		data = data["nearby"].(map[string]interface{})
	}

	if !utils.Contains(utils.GetKeys(data), "search") {
		err := fmt.Errorf("not found search listings '%v' from '%s' '%v'", data, config.Origin, query)
		log.Error(err)
		return 0, 500, err
	}

	data = data["search"].(map[string]interface{})
	qtd := data["totalCount"]
	return int(qtd.(float64)), 200, nil
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

// func unique(sample []models.Property) []models.Property {
// 	var unique []models.Property
// 	type key struct{ value1, value2 string }

// 	m := make(map[key]int)
// 	for _, v := range sample {
// 		k := key{v.Url, v.BusinessType}
// 		if i, ok := m[k]; ok {
// 			// Overwrite previous value per requirement in
// 			// question to keep last matching value.
// 			unique[i] = v
// 		} else {
// 			// Unique key found. Record position and collect
// 			// in result.
// 			m[k] = len(unique)
// 			unique = append(unique, v)
// 		}
// 	}
// 	return unique
// }
