import requests
from model.extract import extract

X = extract()
X = X.drop(columns=['neighborhood', 'unit_types', 'total'])


r = requests.post('http://localhost:8080/predict',
                  json={'input': X.iloc[0].tolist()})

r.json()
