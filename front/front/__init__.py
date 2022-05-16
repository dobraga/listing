from dash import Dash
from dotenv import load_dotenv
import dash_bootstrap_components as dbc


from front import pages


def create_app() -> Dash:
    load_dotenv()

    dash = Dash(
        __name__,
        external_stylesheets=[dbc.themes.BOOTSTRAP, dbc.icons.FONT_AWESOME],
        suppress_callback_exceptions=True,
    )

    dash = pages.init_app(dash)

    return dash


def create_server():
    return create_app().server
