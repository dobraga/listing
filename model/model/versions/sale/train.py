from model.extract import extract

from pathlib import Path
from pickle import dump
from datetime import datetime
from lightgbm import LGBMRegressor
from category_encoders import TargetEncoder
from sklearn.pipeline import make_pipeline
from sklearn.compose import make_column_transformer

FOLDER = Path(__file__).parent

pipeline = make_column_transformer(
    ('passthrough', ['usable_area', 'bedrooms',
     'bathrooms', 'parking_spaces', 'lat', 'lon']),
    (TargetEncoder(), ['neighborhood', 'unit_types'])
)


X = extract('SALE')
y = X.pop('total')

model = make_pipeline(pipeline, LGBMRegressor()).fit(X, y)

now = datetime.now().strftime('%Y_%m_%d_%H_%M_%S')
model_file = FOLDER / f'model_{now}.pkl'
with open(model_file, 'wb') as f:
    dump(model, f)
