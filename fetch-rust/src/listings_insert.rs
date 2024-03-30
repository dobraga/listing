use crate::listings::ParsedListing;
use serde_json::to_string;
use sqlx::postgres::PgPool;

pub async fn insert_listings(
    pool: &PgPool,
    listings: &Vec<ParsedListing>,
) -> Result<(), sqlx::Error> {
    log::info!("inserting {} listings", listings.len());

    for listing in listings {
        let exists: bool =
            sqlx::query_scalar("SELECT EXISTS(SELECT 1 FROM listings WHERE id = $1)")
                .bind(&listing.id)
                .fetch_one(pool)
                .await?;

        if exists {
            sqlx::query("UPDATE listings
                                 SET updated_at = $1, price = $2, monthly_condo_fee = $3, yearly_iptu = $4, active = true
                               WHERE id = $5")
                .bind(listing.updated_at)
                .bind(listing.price)
                .bind(listing.monthly_condo_fee)
                .bind(listing.yearly_iptu)
                .bind(&listing.id)
                .execute(pool)
                .await?;
        } else {
            let query = format!("INSERT INTO listings (id, created_at, updated_at, href, title, description, listing_type, business_type, price, monthly_condo_fee, yearly_iptu, city, location_id, neighborhood, pois_list, state, state_acronym, street, street_number, zip_code, zone, lat, lon, unit_types, usable_area, amenities, bathrooms, bedrooms, parking_spaces, suites, construction_status, media, active)
                VALUES ('{}', '{}', '{}', '{}', '{}', '{}', '{}', '{}', {}, {}, {}, '{}', '{}', '{}', ARRAY{}::TEXT[], '{}', '{}', '{}', '{}', '{}', '{}', {}, {}, '{}', {}, ARRAY{}::TEXT[], {}, {}, {}, {}, '{}', ARRAY{}::TEXT[], true)",
                    &listing.id,
                    &listing.created_at,
                    &listing.updated_at,
                    &listing.href,
                    replace(&listing.title),
                    replace(&listing.description),
                    &listing.listing_type,
                    &listing.business_type,
                    &listing.price,
                    &listing.monthly_condo_fee,
                    &listing.yearly_iptu,
                    &listing.city,
                    &listing.location_id,
                    &listing.neighborhood,
                    replace_vec(&listing.pois_list),
                    &listing.state,
                    &listing.state_acronym,
                    replace(&listing.street),
                    &listing.street_number,
                    &listing.zip_code,
                    &listing.zone,
                    &listing.lat,
                    &listing.lon,
                    &listing.unit_types,
                    &listing.usable_area,
                    replace_vec(&listing.amenities),
                    &listing.bathrooms,
                    &listing.bedrooms,
                    &listing.parking_spaces,
                    &listing.suites,
                    &listing.construction_status,
                    replace_vec(&listing.media)
            );


            match sqlx::query(&query).execute(pool).await {
                Err(e) => {
                    log::error!("error {:?} running:\n{}", e, query);
                    return Err(e);
                },
                Ok(_) => ()
            };
        }
    }
    log::info!("inserted {} listings", listings.len());

    Ok(())
}

fn replace_vec(s: &Vec<String>) -> String {
    let strings: Vec<String> = s.iter().map(|f| f.replace('\'', "`")).collect();

    to_string(&strings)
        .expect(format!("cannot transform to json string {:?}", s).as_str())
        .replace('"', "'")
        .replace('\n', " ")
        .replace("<br>", " ")
}

fn replace(s: &String) -> String {
    to_string(&s)
        .expect(format!("cannot transform to json string {}", s).as_str())
        .replace('\n', " ")
        .replace("<br>", " ")
        .replace('\'', "`")
}
