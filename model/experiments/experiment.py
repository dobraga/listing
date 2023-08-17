from sklearn.metrics import mean_absolute_error as mae, mean_absolute_percentage_error as mape, mean_squared_error
from sklearn.compose import ColumnTransformer

from pandas import DataFrame, Series, concat
from typing import Any, Optional, Callable
from sklearn.pipeline import make_pipeline
from sklearn.model_selection import KFold
from dataclasses import dataclass, field
from sklearn.base import RegressorMixin
from numpy import median, sqrt, zeros
from copy import deepcopy
import logging


def _empty_func(x, y):
    return x, y


@dataclass
class Result:
    name: str
    model: RegressorMixin
    principal_metric: float
    metrics: DataFrame = field(repr=False)
    n_splits: int
    preprocess_train_fn: Callable

    def predict(self, X: DataFrame, y: Series = None) -> Series:
        """
        Predicts the target variable using a trained model.

        Parameters:
            X (DataFrame): The input features.
            y (Series, optional): The target variable. Defaults to None.

        Returns:
            Series: The predicted values of the target variable.
        """

        if y is None:
            return self.model.predict(X)

        m = deepcopy(self.model)
        y_pred = zeros(y.shape)

        for train_index, test_index in KFold(n_splits=self.n_splits, shuffle=True, random_state=42).split(X, y):
            X_train, y_train = self.preprocess_train_fn(
                X.iloc[train_index, :], y[train_index])
            X_test = X.iloc[test_index, :]

            y_pred[test_index] = m.fit(X_train, y_train).predict(X_test)

        return y_pred


@dataclass
class Experiment:
    X: DataFrame = field(repr=False)
    y: Series = field(repr=False)
    n_splits: int
    reduce_metric: Callable = field(
        repr=False, init=False, default_factory=lambda: median)
    metrics: list[Callable] = field(
        repr=True, init=False,
        default_factory=lambda: [mae, rmse, mape])
    obs_metrics: DataFrame = field(
        repr=False, init=False, default_factory=lambda: DataFrame())

    def run(self, name: str, model: RegressorMixin,
            hyperparameters: dict[str, Any] = {},
            preprocessor: Optional[ColumnTransformer] = None,

            X: DataFrame = None,
            y: Series = None,

            preprocess_train_fn: Callable = _empty_func,
            ):
        """
        Runs an experiment using the given model and hyperparameters.

        Parameters:
            name (str): The name of the experiment.
            model (RegressorMixin or callable): The regression model to use. If callable, it is instantiated with the hyperparameters.
            hyperparameters (dict[str, Any], optional): The hyperparameters for the model. Defaults to an empty dictionary.
            preprocessor (ColumnTransformer, optional): The preprocessor to apply to the data before fitting the model. Defaults to None.
            X (DataFrame, optional): The input features for the experiment. If not provided, uses self.X. Defaults to None.
            y (Series, optional): The target variable for the experiment. If not provided, uses self.y. Defaults to None.
            preprocess_train_fn (Callable, optional): The function to preprocess the training data before fitting the model. Defaults to lambda x, y: (x, y).

        Returns:
            Result: The result of the experiment, including the name, fitted model, reduced metric value, and all metrics.
        """

        _metrics: list[dict[str, Any]] = []
        _principal_metrics: list[float] = []

        X = self.X if X is None else X
        y = self.y if y is None else y

        m = model(**hyperparameters) if callable(model) else model
        if preprocessor is not None:
            m = make_pipeline(preprocessor, m)
        LOG.debug(f'Running experiment "{name}":\n"{m}"')

        for split, (train_index, test_index) in enumerate(KFold(n_splits=self.n_splits, shuffle=True, random_state=42).split(X, y)):
            X_train, y_train = preprocess_train_fn(
                X.iloc[train_index, :], y[train_index])
            X_test, y_test = X.iloc[test_index, :], y[test_index]

            m.fit(X_train, y_train)

            y_train_pred = m.predict(X_train)
            y_test_pred = m.predict(X_test)

            for i, metric in enumerate(self.metrics):
                train_metric = metric(y_train, y_train_pred)
                test_metric = metric(y_test, y_test_pred)

                if i == 0:
                    _principal_metrics.append(test_metric)

                mname = metric.__name__
                _metrics.append(
                    {'split_name': 'train', 'split': split,
                     'metric': mname, 'value': train_metric})
                _metrics.append(
                    {'split_name': 'test', 'split': split,
                     'metric': mname, 'value': test_metric})

        reduced = self.reduce_metric(_principal_metrics)
        all_metrics = DataFrame(_metrics)

        self.obs_metrics = concat([self.obs_metrics,
                                   all_metrics.assign(name=name)],
                                  ignore_index=True, axis=0)

        exp = Result(name, m, reduced, all_metrics,
                     self.n_splits, preprocess_train_fn)
        LOG.info(exp)
        return exp


LOG = logging.getLogger(__name__)


def rmse(y_true, y_pred):
    """
    Calculate the root mean squared error (RMSE) between the true values and the predicted values.

    Args:
        y_true (array-like): The true values.
        y_pred (array-like): The predicted values.

    Returns:
        float: The RMSE value.
    """
    return sqrt(mean_squared_error(y_true, y_pred))


mape.__name__ = 'mape'
mae.__name__ = 'mae'
