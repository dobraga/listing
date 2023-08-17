from lightgbm import LGBMRegressor
from pathlib import Path
import m2cgen as m2c
import pickle

from model.extract import extract

X = extract()
X = X.drop(columns=['neighborhood', 'unit_types'])
y = X.pop('total')

model = LGBMRegressor().fit(X, y)
model.predict(X)

code = m2c.export_to_go(model, function_name='Predict')

Path('model/experiments/api/go/model/model.go').write_text('package model\n\n'+code)
Path('model/experiments/api/model.pkl').write_bytes(pickle.dumps(model))
