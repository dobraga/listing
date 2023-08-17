from sqlalchemy import create_engine
import pandas as pd

from model.settings import init


def extract(business_type='RENTAL'):
    """
    Extracts data from a PostgreSQL database based on the specified business type.

    :param business_type: The type of business to extract data for. Defaults to 'RENTAL'.
    :type business_type: str

    :return: A pandas DataFrame containing the extracted data.
    :type: pandas.DataFrame
    """
    settings = init()
    engine = create_engine(
        f"postgresql://{settings['POSTGRES_USER']}:{settings['POSTGRES_PASSWORD']}@localhost:{settings['POSTGRES_PORT']}/{settings['POSTGRES_DB']}")

    return pd.read_sql_query(
        QUERY.format(business_type), engine).set_index(["title", "url", "origin"])


QUERY = """
SELECT title, url, origin, neighborhood, usable_area,
       CASE WHEN unit_types IN ('APARTMENT', 'FLAT', 'KITNET', 'LOFT') THEN 'APARTMENT'
            WHEN unit_types IN ('HOME', 'PENTHOUSE', 'CONDOMINIUM', 'RESIDENTIAL_ALLOTMENT_LAND', 'VILLAGE_HOUSE') THEN 'HOME'
            ELSE 'OTHERS'
        END AS unit_types,
       floors, bedrooms, bathrooms, suites, parking_spaces,
       amenities, lat, lon,
       price + condo_fee AS total,
       created_date
  FROM properties
 WHERE business_type='{}'
   AND created_date >= '2021-01-01'
"""
