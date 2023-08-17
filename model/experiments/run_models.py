from lightgbm import LGBMRegressor
from sklearn.tree import DecisionTreeRegressor
from xgboost import XGBRegressor, XGBRFRegressor
from sklearn.linear_model import LinearRegression, BayesianRidge, Lars, ElasticNet
from sklearn.ensemble import RandomForestRegressor, HistGradientBoostingRegressor, GradientBoostingRegressor, ExtraTreesRegressor

from sklearn.compose import ColumnTransformer
from category_encoders import TargetEncoder
from sklearn.impute import SimpleImputer
from sklearn.svm import SVR, LinearSVR
import pandas as pd
import logging

from model.extract import extract
from experiments.experiment import Experiment


def main():
    X = extract('RENTAL')
    y = X.pop('total')

    exp = Experiment(X, y, 10)

    for model in MODELS:
        for prep_nm, prep in PREPROCESS_TRAIN_FNS.items():
            exp.run(model.__name__ + ' | ' + prep_nm, model,
                    preprocessor=PREPROCESS, preprocess_train_fn=prep)

    metrics = exp.obs_metrics.query('split_name == "test"').groupby(
        ['name', 'metric']).value.median().unstack()

    LOG.info(metrics.sort_values('rmse').head(5))
    LOG.info(metrics.sort_values('mae').head(5))
    LOG.info(metrics.sort_values('mape').head(5))


def clip_target(X, y, quantiles=[.05, .95]):
    qinf = y.quantile(quantiles[0])
    qsup = y.quantile(quantiles[1])
    return X, y.clip(lower=qinf, upper=qsup)


PREPROCESS = ColumnTransformer([
    ('num',
     SimpleImputer(strategy='constant', fill_value=-1),
     ['usable_area', 'floors', 'bedrooms', 'bathrooms', 'suites', 'parking_spaces']),
    ('cat', TargetEncoder(), ['neighborhood', 'unit_types'])
])

PREPROCESS_TRAIN_FNS = {
    'without': lambda x, y: (x, y),
    'clip_01_99': lambda x, y: clip_target(x, y, [0.01, 0.99]),
    'clip_05_95': lambda x, y: clip_target(x, y, [0.05, 0.95]),
    'clip_10_90': lambda x, y: clip_target(x, y, [0.10, 0.90]),
    'clip_25_75': lambda x, y: clip_target(x, y, [0.25, 0.75]),
}

MODELS = [
    LGBMRegressor, XGBRegressor, XGBRFRegressor,
    DecisionTreeRegressor, SVR, LinearSVR,
    BayesianRidge, Lars, ElasticNet, LinearRegression,
    RandomForestRegressor, HistGradientBoostingRegressor,
    GradientBoostingRegressor, ExtraTreesRegressor
]


LOG = logging.getLogger(__name__)

if __name__ == '__main__':
    from rich.logging import RichHandler
    import warnings
    import logging

    logging.basicConfig(
        level="DEBUG",
        format="%(message)s",
        datefmt="[%X]",
        handlers=[RichHandler(rich_tracebacks=True)]
    )

    with warnings.catch_warnings():
        warnings.simplefilter("ignore")
        main()
