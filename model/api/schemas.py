from pydantic import BaseModel
from typing import List


class Input(BaseModel):
    usable_area: float
    lon: float
    lat: float
    parking_spaces: int
    bathrooms: int
    bedrooms: int
    unit_types: str
    neighborhood: str
    condo_fee: float
    price: float


class Inputs(BaseModel):
    data: List[Input]
