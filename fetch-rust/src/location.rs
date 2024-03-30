use crate::{cfg::Config, request::request};

use serde::Deserialize;
use serde_json;
use std::collections::HashMap;
use std::error::Error;

pub async fn list_locations(
    cfg: &Config,
    local: &str,
    origin: &str,
) -> Result<Vec<Address>, Box<dyn Error>> {
    let site = cfg.sites.get(origin).unwrap();

    log::info!("search '{}' in {}", local, site.portal);

    let mut query = HashMap::new();
    query.insert("q", local.to_string());
    query.insert("portal", site.portal.clone());
    query.insert("size", "6".to_string());
    query.insert("fields", "neighborhood,city,street".to_string());
    query.insert("includeFields", "address.city,address.zone,address.state,address.neighborhood,address.stateAcronym,address.street,address.locationId,address.point".to_string());

    let body: String = match request(site.api.as_str(), site.locations.as_str(), &query).await {
        Ok(body) => body,
        Err(err) => return Err(err),
    };

    let location: LocationsResult = match serde_json::from_str(&body) {
        Ok(location) => location,
        Err(err) => return Err(err.into()),
    };

    let addresses = location
        .street
        .result
        .locations
        .iter()
        .map(|l| l.address.clone())
        .collect();
    log::debug!("searching '{}': {:?}", local, addresses);

    Ok(addresses)
}

#[derive(Deserialize)]
struct LocationsResult {
    street: LocationInfo,
    // neighborhood: LocationInfo,
    // city: LocationInfo,
}

#[derive(Deserialize)]
struct LocationInfo {
    // time: i32,
    // maxScore: f32,
    // #[serde(rename = "totalCount")]
    // total_count: i32,
    result: LocationResult,
}

#[derive(Deserialize)]
pub struct LocationResult {
    pub locations: Vec<Location>,
}

#[derive(Deserialize)]
pub struct Location {
    pub address: Address,
}

#[derive(Debug, Deserialize, Clone)]
pub struct Address {
    pub city: String,
    pub zone: String,
    pub street: String,
    #[serde(rename = "locationId")]
    pub location_id: String,
    #[serde(rename = "stateAcronym")]
    pub state_acronym: String,
    pub state: String,
    pub neighborhood: String,
    pub point: Point,
}

#[derive(Debug, Deserialize, Clone)]
pub struct Point {
    pub lon: f32,
    // pub source: String,
    pub lat: f32,
}
