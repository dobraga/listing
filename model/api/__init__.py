from joblib import load
from typing import List
from pathlib import Path
from fastapi import FastAPI
from pandas import DataFrame
from sklearn.base import RegressorMixin

from .schemas import Inputs

MODEL_FOLDER = Path(__file__).parent.parent / 'model'

app = FastAPI()

model_rental = load(
    MODEL_FOLDER/'versions/rental/model_2023_08_15_20_28_01.pkl')
model_sale = load(
    MODEL_FOLDER/'versions/sale/model_2023_08_15_20_27_55.pkl')


@app.get("/health")
async def health():
    return {"status": "ok"}


@app.post("/sale")
async def predict_sale(data: Inputs) -> List[float]:
    data = DataFrame(data.dict()['data'])
    return model_sale.predict(data).tolist()


@app.post("/rental")
async def predict_rental(data: Inputs) -> List[float]:
    data = DataFrame(data.dict()['data'])
    return model_rental.predict(data).tolist()
