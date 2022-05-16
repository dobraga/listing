from dash import dcc, html
import dash_bootstrap_components as dbc

layout = dbc.Navbar(
    html.Div(
        [
            dbc.Button(
                html.I(className="fas fa-bars"),
                color="transparent",
                id="btn_sidebar",
            ),
            html.A(
                html.Span("Imóveis", className="navbar-brand mb-0 h1"),
                href="/",
                style={"color": "inherit"},
            ),
        ]
    )
)
