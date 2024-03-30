use dotenv;
use serde_derive::Deserialize;
use std::collections::HashMap;
use std::env;
use std::fs;
use std::path::Path;
use std::process::exit;

pub fn read_config(filename: &str) -> Config {
    log::debug!("Reading .env");
    let path = Path::new(filename).parent().unwrap().with_file_name(".env");

    match dotenv::from_path(path) {
        Ok(_) => (),
        Err(e) => {
            log::error!("Could not load environment variables: {}", e);
            exit(1);
        }
    }

    log::debug!("Reading config file: {}", filename);
    let contents = match fs::read_to_string(filename) {
        Ok(c) => c,
        Err(e) => {
            log::error!("Could not read file '{}': {}", filename, e);
            exit(1);
        }
    };

    let cfg: Config = match toml::from_str(&contents) {
        Ok(d) => d,
        Err(e) => {
            log::error!("Could not parse file '{}': {}", filename, e);
            exit(1);
        }
    };

    cfg
}

#[derive(Debug, Deserialize)]
pub struct Config {
    pub sites: HashMap<String, Site>,
    pub metro_trem: HashMap<String, Vec<String>>,
    #[serde(default = "default_uri")]
    pub uri: String,
    #[serde(default = "default_backend_host")]
    pub backend_host: String,
    #[serde(default = "default_model_host")]
    pub model_host: String,
    #[serde(default = "default_force_update")]
    pub force_update: bool,
    #[serde(default = "default_max_pages")]
    pub max_pages: u16,
}

fn default_backend_host() -> String {
    let envvar = env::var("ENV").expect("ENV must be set");
    let var_name = envvar + "_BACKEND_HOST";
    env::var(&var_name).expect(format!("{} must be set", &var_name).as_str())
}

fn default_model_host() -> String {
    let envvar = env::var("ENV").expect("ENV must be set");
    let var_name = envvar + "_BACKEND_HOST";
    env::var(&var_name).expect(format!("{} must be set", &var_name).as_str())
}

fn default_force_update() -> bool {
    let envvar = env::var("ENV").expect("ENV must be set");
    let var_name = envvar + "_force_update";
    let value = env::var(&var_name).expect(format!("{} must be set", &var_name).as_str());

    match value.to_lowercase().as_str() {
        "true" | "1" | "yes" => true,
        _ => false,
    }
}

fn default_max_pages() -> u16 {
    let envvar = env::var("ENV").expect("ENV must be set");
    let var_name = envvar + "_max_pages";
    let value: u16 = env::var(&var_name)
        .expect(format!("{} must be set", &var_name).as_str())
        .parse()
        .unwrap();
    value
}

fn default_uri() -> String {
    let envvar = env::var("ENV").expect("ENV must be set");
    let ps_host = envvar + "_POSTGRES_HOST";

    format!(
        "postgres://{}:{}@{}:{}",
        env::var("POSTGRES_USER").expect("POSTGRES_USER must be set"),
        env::var("POSTGRES_PASSWORD").expect("POSTGRES_PASSWORD must be set"),
        env::var(&ps_host).expect(format!("{} must be set", &ps_host).as_str()),
        env::var("POSTGRES_PORT").expect("POSTGRES_PORT must be set"),
    )
}

#[derive(Debug, Deserialize)]
pub struct Site {
    pub site: String,
    pub api: String,
    pub portal: String,
    pub locations: String,
    pub listings: String,
}
