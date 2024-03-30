use crate::listings_type::Type;

use chrono::NaiveDateTime;
use sqlx::postgres::PgPool;

pub async fn last_update(pool: &PgPool, location_id: &str, tp_listings: &Type) -> Option<i64> {
    let last_update: Result<Option<NaiveDateTime>, sqlx::Error> = sqlx::query_scalar(
        "
          SELECT MAX(createdat)
            FROM listings
           WHERE location_id = $1 AND listing_type = $2 AND business_type = $3",
    )
    .bind(location_id)
    .bind(tp_listings.listing_type.to_string())
    .bind(tp_listings.business_type.to_string())
    .fetch_optional(pool)
    .await;

    match last_update {
        Err(e) => {
            log::warn!("[{location_id}] error searching last update: {:?}", e);
            return None;
        }

        Ok(Some(last_update)) => {
            let dt_diff = chrono::Utc::now()
                .signed_duration_since(
                    chrono::DateTime::from_timestamp(last_update.and_utc().timestamp(), 0).unwrap(),
                )
                .num_hours();

            log::info!("[{location_id}] last update {last_update}, {dt_diff} hour(s) ago");
            return Some(dt_diff);
        }

        Ok(None) => {
            log::info!("[{location_id}] not found last update");
            return None;
        }
    };
}

pub async fn reset_active(pool: &PgPool, location_id: &str, tp_listings: &Type) {
    log::info!("[{}] reset active", location_id);
    sqlx::query(
        "
          UPDATE listings
             SET active = false
           WHERE location_id = $1 AND listing_type = $2 AND business_type = $3",
    )
    .bind(location_id)
    .bind(tp_listings.listing_type.to_string())
    .bind(tp_listings.business_type.to_string())
    .execute(pool)
    .await
    .expect("Could not reset active");
}
