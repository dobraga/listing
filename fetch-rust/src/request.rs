use reqwest::header::{HeaderMap, HeaderValue};
use reqwest::{Client, StatusCode};
use std::collections::HashMap;
use std::error::Error;
use url::Url;

pub async fn request(
    url: &str,
    path: &str,
    query: &HashMap<&str, String>,
) -> Result<String, Box<dyn Error>> {
    let client = Client::new();

    let mut headers = HeaderMap::new();
    headers.insert("User-Agent", HeaderValue::from_static("Mozilla/5.0"));
    headers.insert("Accept", HeaderValue::from_static("application/json"));
    headers.insert("x-domain", HeaderValue::from_static(".zapimoveis.com.br"));


    let url = format!("{}{}", url, path);
    let mut url = Url::parse(url.as_str()).expect("Invalid URL");
    for (param, value) in query {
        url.query_pairs_mut().append_pair(param, value);
    }
    let url = url.as_str();
    log::debug!("request url '{}'", url);

    let response = client.get(url).headers(headers).send().await?;

    if response.status() == StatusCode::OK {
        let body = response.text().await?;
        Ok(body)
    } else {
        let url = response.url().to_string();
        let err_msg = format!("Request to '{}' failed with status code: {}", url, response.status());
        Err(err_msg.into())
    }
}
