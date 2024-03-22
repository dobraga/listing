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

pub mod cfg;
pub mod request;
pub mod location;
pub mod properties;

use crate::cfg::read_config;
use crate::location::list_locations;

#[tokio::main]
async fn main() {
    let cfg = read_config("../settings.toml");

    let locations =  list_locations(cfg, "vila valqueire", "vivareal").await.unwrap();
    println!("{:?}", locations.locations[0].address);

    let tp_properties = properties::TypeProperties {
        listing_type: properties::ListingType::USED,
        business_type: properties::BusinessType::RENT,
    };

    let props = properties::fetch_properties(locations.locations[0].address, tp_properties);

}

