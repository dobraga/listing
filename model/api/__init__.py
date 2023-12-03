from joblib import load
import lightgbm as lgb
from typing import List
from pathlib import Path
from pandas import DataFrame
from .schemas import Inputs
from fastapi import FastAPI, HTTPException
from sklearn.model_selection import cross_val_predict
from sklearn.preprocessing import OneHotEncoder
from sklearn.compose import ColumnTransformer
from sklearn.pipeline import Pipeline
from logging import getLogger

app = FastAPI()

MODEL_FOLDER = Path(__file__).parent.parent / 'model'
model_rental_path = MODEL_FOLDER/'versions/rental/model_2023_08_15_20_28_01.pkl'
model_rental = None
if model_rental_path.exists():
    model_rental = load(model_rental_path)

model_sale_path = MODEL_FOLDER/'versions/sale/model_2023_08_15_20_27_55.pkl'
model_sale = None
if model_sale_path.exists():
    model_sale = load(model_sale_path)

ct = ColumnTransformer(
    transformers=[
        ("numerical", "passthrough", [
         "usable_area", "lon", "lat", "parking_spaces", "bathrooms", "bedrooms"]),
        ("categorical", OneHotEncoder(handle_unknown="ignore", min_frequency=0.1),
         ["unit_types", "neighborhood"])
    ]
)
model = Pipeline([("preprocess", ct), ("model", lgb.LGBMRegressor())])


@app.get("/health")
async def health():
    return {"status": "ok"}


@app.post("/sale")
async def predict_sale(data: Inputs) -> List[float]:
    if not model_sale:
        raise HTTPException(status_code=404, detail="Sale model not found")

    data = DataFrame(data.dict()['data'])
    return model_sale.predict(data).tolist()


@app.post("/rental")
async def predict_rental(data: Inputs) -> List[float]:
    if not model_sale:
        raise HTTPException(status_code=404, detail="Rental model not found")

    data = DataFrame(data.dict()['data'])
    return model_rental.predict(data).tolist()


@app.post("/predict")
async def predict_sale(data: Inputs) -> List[float]:
    data = DataFrame(data.dict()['data'])
    y = data.price + data.condo_fee.fillna(0)
    return cross_val_predict(model, data, y, cv=5).tolist()

LOG = getLogger(__name__)
