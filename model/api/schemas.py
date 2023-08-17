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

class Inputs(BaseModel):
    data: List[Input]
