use crate::listings_insert;
use crate::listings_type::Type;
use crate::listings_update;
use crate::{cfg::Config, location::Address as LocationAddress, request::request};

use futures::future::join_all;
use rand::Rng;
use serde::{Deserialize, Serialize};
use sqlx::postgres::PgPool;
use sqlx::FromRow;
use std::collections::HashMap;
use std::error::Error;
use std::time::Duration;

pub async fn fetch_listings(
    pool: &PgPool,
    cfg: &Config,
    address: &LocationAddress,
    type_properties: &Type,
) -> Result<Vec<ParsedListing>, Box<dyn Error>> {
    let name = format!(
        "{} | {} | {}",
        type_properties.listing_type.to_string(),
        type_properties.business_type.to_string(),
        address.location_id
    );
    let mut listings: Vec<ParsedListing> = Vec::new();

    let result = listings_update::last_update(pool, &address.location_id, type_properties).await;
    if let Some(dt_diff) = result {
        if dt_diff < 6 && !cfg.force_update {
            log::info!("[{}] skipping fetch", name);
            return Ok(listings);
        }
        listings_update::reset_active(pool, &address.location_id, type_properties).await;
    } else {
        log::debug!("[{}] last update not found", name)
    }

    log::info!("[{}] fetching", name);
    let mut futures = Vec::new();

    for (site, _) in &cfg.sites {
        let fut = fetch_listings_one_origin(cfg, address, type_properties, site);
        futures.push(fut);
    }

    for result in join_all(futures).await {
        let new_listings: Vec<ParsedListing> = result?;
        listings.extend(new_listings);
    }

    listings_insert::insert_listings(pool, &listings)
        .await
        .expect("Could not insert listings");

    Ok(listings)
}

pub async fn fetch_listings_one_origin(
    cfg: &Config,
    address: &LocationAddress,
    type_properties: &Type,
    origin: &str,
) -> Result<Vec<ParsedListing>, Box<dyn Error>> {
    let mut rng = rand::thread_rng();
    let mut listings: Vec<ParsedListing> = Vec::new();

    log::info!("[{origin}] fetch listings");
    let page_size: u16 = 24;
    let query = create_query(address, type_properties, page_size);

    let (qtd, new_listings) = fetch_listings_page(cfg, &query, origin).await.unwrap();
    listings.extend(new_listings);
    let qtd_pages: u16 = qtd / page_size;
    let max_pages = cfg.max_pages.min(qtd_pages);
    log::info!("[{origin}] have {qtd} listings in {qtd_pages} pages, getting {max_pages}");

    for page in 1..=max_pages {
        log::info!("[{origin}] requesting page {page}");
        let mut new_query = query.clone();
        new_query.insert("from", (page * page_size).to_string());
        let (_, new_listings) = fetch_listings_page(cfg, &new_query, origin).await.unwrap();
        listings.extend(new_listings);

        if page != max_pages {
            let sleep_duration = Duration::from_millis(rng.gen_range(500..=4000));
            log::debug!("sleeping: {:?}", sleep_duration);
            tokio::time::sleep(sleep_duration).await;
        }
    }

    Ok(listings)
}

async fn fetch_listings_page(
    cfg: &Config,
    query: &HashMap<&str, String>,
    origin: &str,
) -> Result<(u16, Vec<ParsedListing>), Box<dyn Error>> {
    let site = cfg.sites.get(origin).unwrap();

    let body: String = match request(site.api.as_str(), site.listings.as_str(), query).await {
        Ok(body) => body,
        Err(err) => {
            log::error!("Could not request {:#?}", err);
            return Err(err);
        }
    };

    // std::fs::write("body.json", &body).unwrap();
    // let jd = &mut serde_json::Deserializer::from_str(body.as_str());
    // let result: Result<Root, _> = serde_path_to_error::deserialize(jd);

    // match result {
    //     Ok(_) => continue,
    //     Err(err) => {
    //         let path = err.path();
    //         log::error!("Could not parse body at {}", path);
    //     }
    // }

    let root: Root = match serde_json::from_str(&body) {
        Ok(root) => root,
        Err(err) => {
            log::error!("Could not parse body {:#?}", err);
            return Err(err.into());
        }
    };
    let qtd = root.search.total_count;

    let listings = parse_listings(origin, root.search.result.listings);

    Ok((qtd, listings))
}

fn create_query<'a>(
    address: &'a LocationAddress,
    type_properties: &'a Type,
    page_size: u16,
) -> HashMap<&'a str, String> {
    let mut query = HashMap::new();

    query.insert("addressNeighborhood", address.neighborhood.clone());
    query.insert("addressLocationId", address.location_id.clone());
    query.insert("addressState", address.state.clone());
    query.insert("addressCity", address.city.clone());
    query.insert("addressZone", address.zone.clone());

    query.insert("listingType", type_properties.listing_type.to_string());
    query.insert("businessType", type_properties.business_type.to_string());

    query.insert("usageTypes", "RESIDENTIAL".to_string());
    query.insert("categoryPage", "RESULT".to_string());
    query.insert("size", page_size.to_string());
    query.insert("from", "0".to_string());

    log::debug!("query: {:?}", query);

    query.insert("includeFields", "search(result(listings(listing(displayAddressType,amenities,usableAreas,conpub structionStatus,listingType,description,title,createdAt,floors,unitTypes,propertyType,id,parkingSpaces,updatedAt,address,suites,externalId,bathrooms,bedrooms,pricingInfos,status),medias,link)),totalCount),expansion(search(result(listings(listing(displayAddressType,amenities,usableAreas,conpub structionStatus,listingType,description,title,createdAt,floors,unitTypes,propertyType,id,parkingSpaces,updatedAt,address,suites,externalId,bathrooms,bedrooms,pricingInfos,status),medias,link)),totalCount)),nearby(search(result(listings(listing(displayAddressType,amenities,usableAreas,conpub structionStatus,listingType,description,title,createdAt,floors,unitTypes,propertyType,id,parkingSpaces,updatedAt,address,suites,externalId,bathrooms,bedrooms,pricingInfos,status),medias,link)),totalCount)),page,developments(search(result(listings(listing(displayAddressType,amenities,usableAreas,conpub structionStatus,listingType,description,title,createdAt,floors,unitTypes,propertyType,id,parkingSpaces,updatedAt,address,suites,externalId,bathrooms,bedrooms,pricingInfos,status),medias,link)),totalCount))".to_string());

    query
}

fn parse_listings(origin: &str, listings: Vec<ListingItem>) -> Vec<ParsedListing> {
    let mut parsed_listings: Vec<ParsedListing> = Vec::new();

    for listing in listings {
        let href = format!("www.{}.com.br{}", origin, listing.link.href);
        let point = listing.listing.address.point;
        let lat = point.lat;
        let lon = point.lon;

        if listing.listing.unit_types.len() > 1 {
            log::warn!(
                "listing with more than 1 unit type: {:?} in '{}'",
                listing.listing.unit_types,
                href
            );
        }
        if listing.listing.usable_areas.len() > 1 {
            log::warn!(
                "listing with more than 1 usable areas: {:?} in '{}'",
                listing.listing.usable_areas,
                href
            );
        }
        if listing.listing.bathrooms.len() > 1 {
            log::warn!(
                "listing with more than 1 bathrooms: {:?} in '{}'",
                listing.listing.bathrooms,
                href
            );
        }
        if listing.listing.bedrooms.len() > 1 {
            log::warn!(
                "listing with more than 1 bedrooms: {:?} in '{}'",
                listing.listing.bedrooms,
                href
            );
        }
        if listing.listing.parking_spaces.len() > 1 {
            log::warn!(
                "listing with more than 1 parking_spaces: {:?} in '{}'",
                listing.listing.parking_spaces,
                href
            );
        }
        if listing.listing.suites.len() > 1 {
            log::warn!(
                "listing with more than 1 suites: {:?} in '{}'",
                listing.listing.suites,
                href
            );
        }

        let usable_area: i32 = listing
            .listing
            .usable_areas
            .first()
            .unwrap_or(&"0".to_string())
            .parse()
            .unwrap();

        let media: Vec<String> = listing
            .medias
            .iter()
            .filter(|m| m.r#type == "IMAGE")
            .map(|m| m.url.clone())
            .collect();

        for l in listing.listing.pricing_infos {
            let price: i64 = l.price.parse().unwrap();
            let monthly_condo_fee: i32 = l.monthly_condo_fee.parse().unwrap();
            let yearly_iptu: i32 = l.yearly_iptu.parse().unwrap();

            parsed_listings.push(ParsedListing {
                id: listing.listing.id.clone(),
                created_at: listing.listing.created_at,
                updated_at: listing.listing.updated_at,
                href: href.clone(),
                title: listing.listing.title.clone(),
                description: listing.listing.description.clone(),
                listing_type: listing.listing.listing_type.clone(),
                business_type: l.business_type.clone(),
                price: price,
                monthly_condo_fee: monthly_condo_fee,
                yearly_iptu: yearly_iptu,
                city: listing.listing.address.city.clone(),
                location_id: listing.listing.address.location_id.clone(),
                neighborhood: listing.listing.address.neighborhood.clone(),
                pois_list: listing.listing.address.pois_list.clone(),
                state: listing.listing.address.state.clone(),
                state_acronym: listing.listing.address.state_acronym.clone(),
                street: listing.listing.address.street.clone(),
                street_number: listing.listing.address.street_number.clone(),
                zip_code: listing.listing.address.zip_code.clone(),
                zone: listing.listing.address.zone.clone(),
                lat: lat,
                lon: lon,
                unit_types: listing.listing.unit_types[0].clone(),
                usable_area: usable_area,
                amenities: listing.listing.amenities.clone(),
                bathrooms: *listing.listing.bathrooms.first().unwrap_or(&0),
                bedrooms: *listing.listing.bedrooms.first().unwrap_or(&0),
                parking_spaces: *listing.listing.parking_spaces.first().unwrap_or(&0),
                suites: *listing.listing.suites.first().unwrap_or(&0),
                construction_status: listing.listing.construction_status.clone(),
                media: media.clone(),
            })
        }
    }

    parsed_listings
}

#[derive(Debug, Serialize, FromRow)]
pub struct ParsedListing {
    pub id: String,

    pub created_at: chrono::DateTime<chrono::Utc>,
    pub updated_at: chrono::DateTime<chrono::Utc>,

    pub href: String,
    pub title: String,
    pub description: String,

    pub listing_type: String,
    pub business_type: String,
    pub price: i64,
    pub monthly_condo_fee: i32,
    pub yearly_iptu: i32,

    pub city: String,
    pub location_id: String,
    pub neighborhood: String,
    pub pois_list: Vec<String>,
    pub state: String,
    pub state_acronym: String,
    pub street: String,
    pub street_number: String,
    pub zip_code: String,
    pub zone: String,

    pub lat: f32,
    pub lon: f32,

    pub unit_types: String,
    pub usable_area: i32,
    pub amenities: Vec<String>,
    pub bathrooms: i8,
    pub bedrooms: i8,
    pub parking_spaces: i8,
    pub suites: i8,
    pub construction_status: String,

    pub media: Vec<String>,
}

fn value_str_default() -> String {
    "0".to_string()
}

fn str_default() -> String {
    "".to_string()
}

fn point_default() -> Point {
    Point { lat: 0.0, lon: 0.0 }
}

#[derive(Debug, Deserialize)]
struct Point {
    lat: f32,
    lon: f32,
}

#[derive(Debug, Deserialize)]
struct Address {
    city: String,
    #[serde(rename = "locationId")]
    location_id: String,
    // name: String,
    neighborhood: String,
    #[serde(default = "point_default")]
    point: Point,
    #[serde(rename = "poisList")]
    pois_list: Vec<String>,
    state: String,
    #[serde(rename = "stateAcronym")]
    state_acronym: String,
    #[serde(default = "value_str_default")]
    street: String,
    #[serde(rename = "streetNumber", default = "value_str_default")]
    street_number: String,
    #[serde(rename = "zipCode")]
    zip_code: String,
    zone: String,
}

#[derive(Debug, Deserialize)]
struct PricingInfo {
    #[serde(rename = "businessType")]
    business_type: String,
    #[serde(rename = "monthlyCondoFee", default = "value_str_default")]
    monthly_condo_fee: String,
    #[serde(default = "value_str_default")]
    price: String,
    #[serde(rename = "yearlyIptu", default = "value_str_default")]
    yearly_iptu: String,
}

#[derive(Debug, Deserialize)]
struct Listing {
    // #[serde(rename = "externalId")]
    // external_id: String,
    id: String,

    #[serde(rename = "createdAt")]
    created_at: chrono::DateTime<chrono::Utc>,
    #[serde(rename = "updatedAt")]
    updated_at: chrono::DateTime<chrono::Utc>,
    address: Address,
    amenities: Vec<String>,
    bathrooms: Vec<i8>,
    bedrooms: Vec<i8>,
    #[serde(rename = "constructionStatus", default = "str_default")]
    construction_status: String,
    description: String,
    #[serde(rename = "listingType")]
    listing_type: String,
    #[serde(rename = "parkingSpaces")]
    parking_spaces: Vec<i8>,
    #[serde(rename = "pricingInfos")]
    pricing_infos: Vec<PricingInfo>,
    suites: Vec<i8>,
    title: String,
    #[serde(rename = "unitTypes")]
    unit_types: Vec<String>,
    #[serde(rename = "usableAreas")]
    usable_areas: Vec<String>,
}

#[derive(Debug, Deserialize)]
struct Link {
    href: String,
    // name: String,
    // rel: String,
}

#[derive(Debug, Deserialize)]
struct Media {
    r#type: String,
    url: String,
}

#[derive(Debug, Deserialize)]
struct ListingItem {
    link: Link,
    listing: Listing,
    medias: Vec<Media>,
}

#[derive(Debug, Deserialize)]
struct SearchResult {
    listings: Vec<ListingItem>,
}

#[derive(Debug, Deserialize)]
struct Search {
    result: SearchResult,
    #[serde(rename = "totalCount")]
    total_count: u16,
}

#[derive(Debug, Deserialize)]
struct Root {
    search: Search,
}
