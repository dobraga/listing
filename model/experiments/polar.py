import numpy as np
from pandas import DataFrame, Series
from sklearn.base import TransformerMixin, BaseEstimator


def cart_to_pol(data):
    x, y = data['lat'], data['lon']
    complex_format = x + 1j * y
    return {'distance': np.abs(complex_format), 'angle': np.angle(complex_format, deg=True)}


class LatLonPolar(TransformerMixin, BaseEstimator):
    def __init__(self, reference_column='neighborhood'):
        super().__init__()

        self.reference_column = reference_column

        self.distance = f'{reference_column}_distance'
        self.angle = f'{reference_column}_angle'

    def fit(self, X: DataFrame, y: Series = None, **fit_params) -> 'LatLonPolar':

        self.mapping = X.groupby([self.reference_column])[
            ['lat', 'lon']].median()
        self.mapping.columns = ['lat_point', 'lon_point']

        return self

    def transform(self, X: DataFrame) -> DataFrame:
        data = X.set_index(self.reference_column)[['lat', 'lon']]
        data = data.join(self.mapping, how='left')

        data['lat'] = data['lat'] - data['lat_point']
        data['lon'] = data['lon'] - data['lon_point']

        data[[self.distance, self.angle]] = data[
            ['lat', 'lon']].apply(cart_to_pol, axis=1, result_type='expand')

        return data[[self.distance, self.angle]]

    def get_feature_names(self, input_features=None):
        return [self.distance, self.angle]
