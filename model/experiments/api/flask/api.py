from flask import Flask, request
from pathlib import Path
import pickle

model = pickle.loads(Path('../model.pkl').read_bytes())
app = Flask(__name__)


@app.post("/predict")
def predict():
    data = [request.json['input']]
    return model.predict(data).tolist()


if __name__ == "__main__":
    app.run(port=8080)
