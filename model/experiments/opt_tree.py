import logging
import mlflow
import optuna
from mlflow.utils.mlflow_tags import MLFLOW_PARENT_RUN_ID
from mlflow.tracking import MlflowClient

from sklearn.compose import ColumnTransformer
from category_encoders import TargetEncoder
from sklearn.impute import SimpleImputer

from sklearn.ensemble import RandomForestRegressor, ExtraTreesRegressor
from lightgbm import LGBMRegressor
from xgboost import XGBRegressor

from experiments.experiment import Experiment
from model.extract import extract

client = MlflowClient()


def main():
    LOG.debug(f'Tracking URI: {mlflow.get_tracking_uri()}')
    LOG.debug(f'Artifact URI: {mlflow.get_artifact_uri()}')

    experiment, study_run_id = get_study()
    LOG.debug(f'Experiment: {experiment}')
    LOG.debug(f'Study run: {study_run_id}')

    study = optuna.create_study(
        study_name="opt",
        direction="minimize",
        load_if_exists=True,
        storage="sqlite:///opt.db",
    )

    X = extract()
    y = X.pop('total')
    objective = get_objective(X, y, experiment, study_run_id)
    study.optimize(objective, n_trials=1000)


def get_objective(X, y, experiment, study_run_id):
    def objective(trial):
        trial_run = client.create_run(
            experiment_id=experiment, tags={
                MLFLOW_PARENT_RUN_ID: study_run_id}
        ).info.run_id

        name_model = trial.suggest_categorical(
            "model",
            [
                "LGBMRegressor",
                "XGBRegressor",
                "RandomForestRegressor",
                "ExtraTreesRegressor",
            ],
        )

        if name_model == "ExtraTreesRegressor":
            param = {
                "random_state": 20,
                "n_estimators": trial.suggest_int("n_estimators", 10, 500),
                "max_depth": trial.suggest_int("max_depth", 1, 10),
                "max_features": trial.suggest_categorical(
                    "max_features", ["auto", "sqrt", "log2"]
                ),
                "min_samples_leaf": trial.suggest_int("min_samples_leaf", 1, 10),
                "min_samples_split": trial.suggest_int("min_samples_split", 2, 10),
                "criterion": trial.suggest_categorical("criterion", ['friedman_mse', 'squared_error', 'absolute_error', 'poisson']),
                "bootstrap": trial.suggest_categorical("bootstrap", [True, False]),
            }

        elif name_model == "RandomForestRegressor":
            param = {
                "random_state": 20,
                "n_estimators": trial.suggest_int("n_estimators", 10, 500),
                "max_depth": trial.suggest_int("max_depth", 1, 10),
                "max_features": trial.suggest_categorical(
                    "max_features", ["auto", "sqrt", "log2"]
                ),
                "min_samples_leaf": trial.suggest_int("min_samples_leaf", 1, 10),
                "min_samples_split": trial.suggest_int("min_samples_split", 2, 10),
                "criterion": trial.suggest_categorical("criterion", ["mse", "mae"]),
                "bootstrap": trial.suggest_categorical("bootstrap", [True, False]),
            }

        elif name_model == "LGBMRegressor":
            param = {
                "random_state": 20,
                "n_estimators": trial.suggest_int("n_estimators", 10, 500),
                "max_depth": trial.suggest_int("max_depth", 1, 10),
                "colsample_bytree": trial.suggest_float(
                    "colsample_bytree", 0.1, 0.9
                ),
                "subsample": trial.suggest_float("subsample", 0.6, 1.0),
                "num_leaves": trial.suggest_int("num_leaves", 2, 90),
                "min_split_gain": trial.suggest_float("min_split_gain", 0.001, 0.1),
                "reg_alpha": trial.suggest_float("reg_alpha", 0, 1),
                "reg_lambda": trial.suggest_float("reg_lambda", 0, 1),
                "min_child_weight": trial.suggest_int("min_child_weight", 5, 50),
                "learning_rate": trial.suggest_float("learning_rate", 1e-5, 5e-1),
            }

        elif name_model == "XGBRegressor":
            param = {
                "random_state": 20,
                "n_estimators": trial.suggest_int("n_estimators", 10, 500),
                "max_depth": trial.suggest_int("max_depth", 1, 10),
                "colsample_bytree": trial.suggest_float(
                    "colsample_bytree", 0.1, 0.9
                ),
                "subsample": trial.suggest_float("subsample", 0.6, 1.0),
                "reg_alpha": trial.suggest_float("reg_alpha", 0, 1),
                "reg_lambda": trial.suggest_float("reg_lambda", 0, 1),
                "min_child_weight": trial.suggest_int("min_child_weight", 5, 50),
                "learning_rate": trial.suggest_float("learning_rate", 1e-5, 5e-1),
            }

        # Log dos parâmetros
        client.log_param(trial_run, "model", name_model)
        [client.log_param(trial_run, k, v) for k, v in param.items()]

        # Realiza validação cruzada e faz log da métricas
        exp = Experiment.run(X, y, MODELS[name_model], param,
                             preprocessor=PREPROCESS)

        client.log_metric(trial_run, 'mae', exp.mae)
        client.log_metric(trial_run, 'mape', exp.mape)
        client.log_metric(trial_run, 'rmse', exp.rmse)
        client.log_metric(trial_run, 'mse', exp.mse)

        return exp.mape

    return objective


def get_study():
    # https://simonhessner.de/mlflow-optuna-parallel-hyper-parameter-optimization-and-logging/
    client = MlflowClient()
    experiment_name = "opt"

    try:
        experiment = client.create_experiment(experiment_name)
    except:
        experiment = client.get_experiment_by_name(
            experiment_name).experiment_id

    study_run = client.create_run(experiment_id=experiment)
    return experiment, study_run.info.run_id


MODELS = {
    'XGBRegressor': XGBRegressor,
    'LGBMRegressor': LGBMRegressor,
    'RandomForestRegressor': RandomForestRegressor,
    'ExtraTreesRegressor': ExtraTreesRegressor
}

PREPROCESS = ColumnTransformer([
    ('num', SimpleImputer(strategy='constant', fill_value=-1), ['usable_area', 'floors',
     'bedrooms', 'bathrooms', 'suites', 'parking_spaces']),
    ('cat', TargetEncoder(), ['neighborhood', 'unit_types'])
])

LOG = logging.getLogger(__name__)

if __name__ == "__main__":
    from rich.logging import RichHandler
    import warnings
    import logging

    logging.basicConfig(
        level="NOTSET",
        format="%(message)s",
        datefmt="[%X]",
        handlers=[RichHandler(rich_tracebacks=True)]
    )

    with warnings.catch_warnings():
        warnings.simplefilter("ignore")
        main()
