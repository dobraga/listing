from dotenv import dotenv_values
from pathlib import Path
from os import chdir
import logging


def init():
    while not Path('.env').is_file():
        LOG.debug(
            f'"{Path(".env").absolute()}" not exists, return one folder.')
        chdir('..')

    LOG.info(f'"{Path(".env").absolute()}" loading this env file.')

    config = dotenv_values()
    env = config['ENV']

    for var in ENV_VARIABLES:
        config[var] = config[f'{env}_{var}']

    for var in POSTGRES_VARIABLES:
        if var not in config:
            raise KeyError(f'.env file not found "{var}"')

    return config


LOG = logging.getLogger(__name__)


ENV_VARIABLES = ['BACKEND_HOST', 'POSTGRES_HOST',
                 'DEBUG', 'force_update']
POSTGRES_VARIABLES = ['POSTGRES_USER', 'POSTGRES_PASSWORD',
                      'POSTGRES_DB', 'POSTGRES_PORT']
