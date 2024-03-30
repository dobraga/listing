use clap::{command, ArgAction, Parser};
use clap_verbosity_flag::Verbosity;

use crate::listings_type;

#[derive(Parser, Debug)]
#[command(version, about, long_about = "Get and store listings.")]
pub struct Args {
    /// Location
    #[arg(long)]
    pub location: String,

    /// Listing type
    #[arg(short = 'l', long)]
    pub r#type_listing: listings_type::ListingType,

    /// Business type
    #[arg(short = 'b', long)]
    pub r#type_business: listings_type::BusinessType,

    /// Force update
    #[arg(short, long, action=ArgAction::SetTrue)]
    pub force: bool,

    /// Verbose log
    #[command(flatten)]
    pub verbose: Verbosity,
}
