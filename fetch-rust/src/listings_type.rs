#[derive(Debug)]
pub struct Type {
    pub listing_type: ListingType,
    pub business_type: BusinessType,
}

#[derive(Debug, Clone, clap::ValueEnum)]
#[clap(rename_all = "kebab_case")]
pub enum ListingType {
    DEVELOPMENT,
    USED,
}

#[derive(Debug, Clone, clap::ValueEnum)]
#[clap(rename_all = "kebab_case")]
pub enum BusinessType {
    SALE,
    RENT,
}

impl ListingType {
    pub fn to_string(&self) -> String {
        match self {
            ListingType::DEVELOPMENT => "DEVELOPMENT".to_string(),
            ListingType::USED => "USED".to_string(),
        }
    }
}

impl BusinessType {
    pub fn to_string(&self) -> String {
        match self {
            BusinessType::SALE => "SALE".to_string(),
            BusinessType::RENT => "RENTAL".to_string(),
        }
    }
}
