from dash import Dash
import dash_bootstrap_components as dbc


from front import pages
from front.settings import init


def create_app() -> Dash:
    dash = Dash(
        __name__,
        external_stylesheets=[dbc.themes.BOOTSTRAP, dbc.icons.FONT_AWESOME],
        suppress_callback_exceptions=True,
    )
    settings = init()
    dash = pages.init_app(dash, settings)

    return dash


def create_server():
    return create_app().server
