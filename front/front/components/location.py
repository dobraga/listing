from dataclasses import dataclass
from typing import Optional
from requests import get


@dataclass
class Location:
    city: str
    zone: str
    state: str
    locationId: str
    neighborhood: str
    stateAcronym: str
    addressStreet: Optional[str]
    addressPointLat: Optional[float]
    addressPointLon: Optional[float]

    def value(self) -> str:
        if self.addressStreet:
            return self.addressStreet+self.locationId
        return self.locationId

    def label(self) -> str:
        desc = ", ".join([self.stateAcronym, self.city, self.neighborhood])
        if self.addressStreet:
            return f'{self.addressStreet}, {desc}'
        return desc

    def dict(self):
        _dict = self.__dict__.copy()
        _dict['value'] = self.value()
        _dict['label'] = self.label()
        return _dict

    @classmethod
    def parse(cls, url: str) -> list["Location"]:
        request = get(url)
        request.raise_for_status()
        locations = request.json()

        if "adressStreet" in locations[0].keys():
            return [
                cls(
                    city=l["city"],
                    zone=l["zone"],
                    state=l["state"],
                    locationId=l["locationId"],
                    neighborhood=l["neighborhood"],
                    stateAcronym=l["stateAcronym"],
                    addressStreet=l["addressStreet"],
                    addressPointLat=l["addressPointLat"],
                    addressPointLon=l["addressPointLon"],
                )
                for l in locations]

        return [
            cls(
                city=l["city"],
                zone=l["zone"],
                state=l["state"],
                locationId=l["locationId"],
                neighborhood=l["neighborhood"],
                stateAcronym=l["stateAcronym"],
                addressStreet=None,
                addressPointLat=None,
                addressPointLon=None,
            )
            for l in locations]
