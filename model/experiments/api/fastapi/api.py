from pydantic import BaseModel
from fastapi import FastAPI
from pathlib import Path
import pickle

model = pickle.loads(Path('../model.pkl').read_bytes())
app = FastAPI()


class Input(BaseModel):
    input: list[list[float]]


@app.post("/predict")
async def predict(input: Input):
    return model.predict(input.input).tolist()
