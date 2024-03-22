use reqwest::header::{HeaderMap, HeaderValue};
use reqwest::{Client, StatusCode};
use std::collections::HashMap;
use std::error::Error;
use url::Url;


pub async fn request(url: &str, path: &str, query: HashMap<&str, &str>) -> Result<String, Box<dyn Error>> {
    let client = Client::new();

    let mut headers = HeaderMap::new();
    headers.insert("User-Agent", HeaderValue::from_static("Mozilla/5.0"));

    let url = format!("{}{}", url, path);
    let mut url = Url::parse(url.as_str()).expect("Invalid URL");
    for (param, value) in query {
        url.query_pairs_mut().append_pair(param, value);
    }
    let url = url.as_str();

    let response = client
        .get(url)
        .headers(headers)
        .send()
        .await?;

    if response.status() == StatusCode::OK {
        let body = response.text().await?;
        Ok(body)
    } else {
        Err(response.status().to_string().into())
    }
}
