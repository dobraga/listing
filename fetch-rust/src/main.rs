// use axum::{
//     http::StatusCode,
//     routing::{get, post},
//     Json, Router,
// };
// use serde::{Deserialize, Serialize};
// use tokio::net::TcpListener;
// use log::{info};

// #[tokio::main]
// async fn main() {
//     let addr = "0.0.0.0:3000";
//     let app = Router::new()
//         .route("/", get(root))
//         .route("/users", post(create_user));

//     let listener = TcpListener::bind(addr).await.unwrap();
//     info!("Listening on http://{}", addr);
//     axum::serve(listener, app).await.unwrap();
// }

// async fn root() -> &'static str {
//     "Olá mundo!"
// }

// async fn create_user(Json(payload): Json<CreateUser>) -> (StatusCode, Json<User>) {
//     let user = User {
//         id: 1337,
//         username: payload.username,
//     };

//     (StatusCode::CREATED, Json(user))
// }

// #[derive(Deserialize)]
// struct CreateUser {
//     username: String,
// }

// #[derive(Serialize)]
// struct User {
//     id: u64,
//     username: String,
// }

pub mod args;
pub mod cfg;
pub mod db;
pub mod listings;
pub mod listings_insert;
pub mod listings_type;
pub mod listings_update;
pub mod location;
pub mod request;

use crate::args::Args;
use crate::cfg::read_config;
use crate::db::connect_db;
use crate::location::list_locations;
use clap::Parser;

#[tokio::main]
async fn main() {
    let args = Args::parse();

    let level = args
        .verbose
        .log_level_filter()
        .to_level()
        .unwrap_or(log::Level::Info);
    simple_logger::init_with_level(level).unwrap();

    let cfg = read_config("../settings.toml");

    let pool = connect_db(&cfg).await;

    let locations = list_locations(&cfg, &args.location, "vivareal")
        .await
        .unwrap();

    let tp_listings: listings_type::Type = listings_type::Type {
        listing_type: args.type_listing,
        business_type: args.type_business,
    };

    let _vec_listings = listings::fetch_listings(&pool, &cfg, &locations[0], &tp_listings)
        .await
        .expect("Could not fetch listings");
}
