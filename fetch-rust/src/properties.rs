use crate::location::Address;
use std::error::Error;

use std::collections::HashMap;

pub fn fetch_properties(address: Address, type_properties: TypeProperties) {
    // pub fn fetch_properties(address: Address) -> Result<(), Box<dyn Error>> {

    let query = create_query(address, type_properties);
}

pub struct TypeProperties {
    pub listing_type: ListingType,
    pub business_type: BusinessType,
}

pub enum ListingType {
    DEVELOPMENT,
    USED,
}

impl ListingType {
    pub fn to_str(&self) -> &'static str {
        match self {
            ListingType::DEVELOPMENT => "DEVELOPMENT",
            ListingType::USED => "USED",
        }
    }
}

pub enum BusinessType {
    SALE,
    RENT,
}

impl BusinessType {
    pub fn to_str(&self) -> &'static str {
        match self {
            BusinessType::SALE => "SALE",
            BusinessType::RENT => "RENT",
        }
    }
}

fn create_query(
    address: Address,
    type_properties: TypeProperties,
) -> HashMap<&'static str, &'static str> {
    let mut query = HashMap::new();
    query.insert("addressNeighborhood", address.neighborhood);
    query.insert("addressLocationId", address.location_id);
    query.insert("addressState", address.state);
    query.insert("addressCity", address.city);
    query.insert("addressZone", address.zone);

    query.insert("listingType", type_properties.listing_type.to_str());
    query.insert("business", type_properties.business_type.to_str());

    query.insert("usageTypes", "RESIDENTIAL");
    query.insert("categoryPage", "RESULT");
    query.insert("size", "24");
    query.insert("from", "0");

    query
}
