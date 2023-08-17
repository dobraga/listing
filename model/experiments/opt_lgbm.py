from sklearn.compose import ColumnTransformer
from category_encoders import TargetEncoder
from sklearn.impute import SimpleImputer
import lightgbm as lgb
import optuna
import os

from model.extract import extract
from experiments.experiment import Experiment


def objective(trial: optuna.Trial):
    params = {
        "verbosity": -1,
        "boosting_type": "gbdt",
        "random_state": 42,
        "n_jobs": os.cpu_count() - 1,
        "n_estimators": trial.suggest_int("n_estimators", 100, 1000),
        "learning_rate": trial.suggest_float("learning_rate", 1e-8, 1.0),
        "max_depth": trial.suggest_int("max_depth", 1, 30),
        "reg_alpha": trial.suggest_float("reg_alpha", 1e-8, 10.0, log=True),
        "reg_lambda": trial.suggest_float("reg_lambda", 1e-8, 10.0, log=True),
        "num_leaves": trial.suggest_int("num_leaves", 2, 256),
        "colsample_bytree": trial.suggest_float("colsample_bytree", 0.4, 1.0),
        "subsample": trial.suggest_float("subsample", 0.4, 1.0),
        "subsample_freq": trial.suggest_int("subsample_freq", 1, 7),
        "min_child_samples": trial.suggest_int("min_child_samples", 5, 100),
    }

    features = {
        'neighborhood': trial.suggest_int('neighborhood', 0, 1),
        'unit_types': trial.suggest_int('unit_types', 0, 1),
        'usable_area': trial.suggest_int('usable_area', 0, 1),
        'floors': trial.suggest_int('floors', 0, 1),
        'bedrooms': trial.suggest_int('bedrooms', 0, 1),
        'bathrooms': trial.suggest_int('bathrooms', 0, 1),
        'suites': trial.suggest_int('suites', 0, 1),
        'parking_spaces': trial.suggest_int('parking_spaces', 0, 1),
    }

    if sum(features.values()) == 0:
        raise optuna.TrialPruned()

    prep = create_preprocessor(features)
    prep_train = PREPROCESS_TRAIN_FNS[trial.suggest_categorical(
        'preprocess', PREPROCESS_TRAIN_FNS.keys())]

    return exp.run('', lgb.LGBMRegressor, params, preprocessor=prep, preprocess_train_fn=prep_train).principal_metric


def create_preprocessor(features: dict) -> ColumnTransformer:
    """
    Create a preprocessor for the given features.

    This function takes in a dictionary of features and returns a ColumnTransformer object that can be used to preprocess the data. The features dictionary should have the feature names as keys and a boolean value indicating whether the feature should be included or not.

    Parameters:
        - features (dict): A dictionary of features to include in the preprocessor. The keys are the feature names, and the values are boolean flags indicating whether the feature should be included or not.

    Returns:
        - prep (ColumnTransformer): A ColumnTransformer object that can be used to preprocess the data. The preprocessor will apply different transformations to the numerical features and categorical features based on the specified flags in the features dictionary.
    """
    nu_features = [feature for feature,
                   fl in features.items() if feature in NUM_FEATURES and fl]

    ca_features = [feature for feature,
                   fl in features.items() if feature in CAT_FEATURES and fl]

    prep = []

    if nu_features:
        prep.append(('num', SimpleImputer(
            strategy='constant', fill_value=-1), nu_features))

    if ca_features:
        prep.append(('cat', TargetEncoder(), ca_features))

    return ColumnTransformer(prep)


def clip_target(X, y, quantiles=[.05, .95]):
    """
    Clip the target variable `y` within a specified range defined by the lower and upper quantiles.

    Parameters:
        X (pandas.DataFrame): The input features.
        y (pandas.Series): The target variable.
        quantiles (list, optional): The lower and upper quantiles used to define the range. Defaults to [.05, .95].

    Returns:
        tuple: A tuple containing the modified input features `X` and the clipped target variable `y`.
    """
    return X, y.clip(lower=y.quantile(quantiles[0]), upper=y.quantile(quantiles[1]))


CAT_FEATURES = ['neighborhood', 'unit_types']
NUM_FEATURES = ['usable_area', 'floors', 'bedrooms',
                'bathrooms', 'suites', 'parking_spaces']

PREPROCESS_TRAIN_FNS = {
    'without': lambda x, y: (x, y),
    'clip_01_99': lambda x, y: clip_target(x, y, [0.01, 0.99]),
    'clip_05_95': lambda x, y: clip_target(x, y, [0.05, 0.95]),
    'clip_10_90': lambda x, y: clip_target(x, y, [0.10, 0.90]),
    'clip_25_75': lambda x, y: clip_target(x, y, [0.25, 0.75]),
}


if __name__ == "__main__":
    from rich.logging import RichHandler
    import warnings
    import logging

    logging.basicConfig(
        level="DEBUG",
        format="%(message)s",
        datefmt="[%X]",
        handlers=[RichHandler(rich_tracebacks=True)]
    )

    X = extract()
    y = X.pop('total')
    exp = Experiment(X, y, 3)

    with warnings.catch_warnings():
        warnings.simplefilter("ignore")
        study = optuna.create_study(direction="minimize")
        study.optimize(objective, n_trials=100)

    print("Number of finished trials: {}".format(len(study.trials)))

    print("Best trial:")
    trial = study.best_trial

    print("  Value: {}".format(trial.value))

    print("  Params: ")
    for key, value in trial.params.items():
        print("    {}: {}".format(key, value))
