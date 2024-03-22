pub mod properties;

use location::Address;
use std::error::Error;

pub fn fetch_properties(address: Address) -> Result<(), Box<dyn Error>> {

}


fn create_query(address: Address) -> HashMap<&str, &str> {
    let mut query = HashMap::new();
    query.insert("locationId", address.location_id.as_str());
    query
}
