from dash import Dash
from dash import dcc, html

from front.pages import sidebar, navbar, table, map


layout = dcc.Loading(
    id="loading",
    children=html.Div(
        [
            dcc.Store(id="locations"),
            dcc.Store(id="data"),
            dcc.Store(id="filtered_data"),
            dcc.Store(id="side_click"),
            dcc.Location(id="url"),
            navbar.layout,
            sidebar.layout,
            html.Div(id="page-content"),
        ],
    ),
)


def init_app(dash: Dash, settings: dict) -> Dash:
    dash.layout = layout
    dash = sidebar.init_app(dash, settings)
    dash = table.init_app(dash)
    dash = map.init_app(dash)
    return dash
