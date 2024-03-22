use crate::cfg::Config;
use crate::request::request;

use serde_json;
use serde::Deserialize;
use std::collections::HashMap;
use std::error::Error;

pub async fn list_locations(
    cfg: Config,
    local: &str,
    origin: &str,
) -> Result<LocationResult, Box<dyn Error>> {
    let site = cfg.sites.get(origin).unwrap();

    let mut query = HashMap::new();
    query.insert("q", local);
    query.insert("portal", site.portal.as_str());
    query.insert("size", "6");
    query.insert("fields", "neighborhood,city,street");
    query.insert("includeFields", "address.city,address.zone,address.state,address.neighborhood,address.stateAcronym,address.street,address.locationId,address.point");

    let body: String = match request(site.api.as_str(), "/v3/locations", query).await {
        Ok(body) => body,
        Err(err) => return Err(err),
    };

    let location: LocationsResult = match serde_json::from_str(&body) {
        Ok(location) => location,
        Err(err) => return Err(err.into()),
    };

    Ok(location.street.result)
}

#[derive(Debug, Deserialize)]
struct LocationsResult {
    street: LocationInfo,
    // neighborhood: LocationInfo,
    // city: LocationInfo,
}

#[derive(Debug, Deserialize)]
struct LocationInfo {
    // time: i32,
    // maxScore: f32,
    // totalCount: i32,
    result: LocationResult,
}

#[derive(Debug, Deserialize)]
pub struct LocationResult {
    pub locations: Vec<Location>,
}

#[derive(Debug, Deserialize)]
pub struct Location {
    pub address: Address,
}

#[derive(Debug, Deserialize)]
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

#[derive(Debug, Deserialize)]
pub struct Point {
    pub lon: f64,
    // pub source: String,
    pub lat: f64,
}
