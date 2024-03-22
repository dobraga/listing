pub mod cfg {
    pub use crate::cfg::read_config;
    pub struct Config;
    pub struct Site;
}

use serde_derive::Deserialize;
use std::collections::HashMap;
use std::fs;
use std::process::exit;

pub fn read_config(filename: &str) -> Config {
    let contents = match fs::read_to_string(filename) {
        Ok(c) => c,
        Err(e) => {
            eprintln!("Could not read file '{}': {}", filename, e);
            exit(1);
        }
    };

    let cfg: Config = match toml::from_str(&contents) {
        Ok(d) => d,
        Err(e) => {
            eprintln!("Unable to load data from '{}': {}", filename, e);
            exit(1);
        }
    };

    cfg
}

#[derive(Debug, Deserialize)]
pub struct Config {
    pub sites: HashMap<String, Site>,
    pub metro_trem: HashMap<String, Vec<String>>,
}

#[derive(Debug, Deserialize)]
pub struct Site {
    pub site: String,
    pub api: String,
    pub portal: String,
}
