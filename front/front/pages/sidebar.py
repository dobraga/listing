from math import ceil

from dash.dependencies import Input, Output, State
from dash.exceptions import PreventUpdate
import dash_bootstrap_components as dbc
from dash import Dash, dcc, html
from sqlalchemy import create_engine
from logging import getLogger
from requests import get
import pandas as pd

from front.pages import map, table
from front.components.location import Location

depara_tp_contrato = [
    {"label": "Aluguel", "value": "RENTAL"},
    {"label": "Compra", "value": "SALE"},
]
depara_tp_listings = [
    {"label": "Usado", "value": "USED"},
    {"label": "Novo", "value": "DEVELOPMENT"},
]


def min_max(values: pd.Series, plus=0):
    ini, fim = values.min(), values.max()
    ini, fim = int(int(ini / plus)) * plus, int(ceil(fim / plus)) * plus
    return ini, ini, fim, fim, ini, fim


layout = html.Div(
    [
        dbc.Nav(
            [
                dbc.NavLink("Tabela", href="/", active="exact"),
                dbc.NavLink("Mapa", href="/dash/map", active="exact"),
            ],
            vertical=True,
            pills=True,
        ),
        html.Hr(),
        html.Div(
            [
                html.Div(
                    [
                        dbc.Label("Tipo do Imóvel"),
                        dbc.InputGroup(
                            [
                                dbc.Select(
                                    id="business_type",
                                    options=depara_tp_contrato,
                                    value="RENTAL",
                                ),
                                dbc.Select(
                                    id="listing_type",
                                    value="USED",
                                    options=depara_tp_listings,
                                ),
                            ],
                        ),
                    ],
                    className="group_input",
                ),
                html.Div(
                    [
                        dbc.Label("Busca Local"),
                        dbc.InputGroup(
                            [
                                html.Div([
                                    dbc.Label("Tipo da localização"),
                                    dbc.RadioItems(
                                        options=[
                                            {"label": "Rua", "value": "street"},
                                            {"label": "Bairro",
                                                "value": "neighborhood"},
                                        ],
                                        value="neighborhood",
                                        id="type_location",
                                    )
                                ]),
                                dbc.Input(
                                    id="local",
                                    placeholder="Local",
                                ),
                                dbc.Button(
                                    "Procurar", color="primary", id="search_local"
                                ),
                            ]
                        ),
                    ],
                    className="group_input",
                ),
                dbc.Fade(
                    html.Div(
                        [
                            dbc.Select(id="select_local"),
                            dbc.Button(
                                "Procurar Imóveis", color="primary", id="search_imoveis"
                            ),
                        ]
                    ),
                    id="fade_search_imoveis",
                    is_in=False,
                ),
            ]
        ),
        dbc.Fade(
            html.Div(
                [
                    html.Hr(),
                    html.Div(
                        [
                            dbc.Label("Faixa de Preço"),
                            html.Div(
                                [
                                    dcc.Input(
                                        id="preco_min",
                                        type="number",
                                        step=100,
                                    ),
                                    dcc.Input(
                                        id="preco_max",
                                        type="number",
                                        step=100,
                                    ),
                                ]
                            ),
                            html.Br(),
                            dbc.Label("Quantidade de quartos"),
                            html.Div(
                                [
                                    dcc.Input(
                                        id="quarto_min",
                                        type="number",
                                    ),
                                    dcc.Input(
                                        id="quarto_max",
                                        type="number",
                                    ),
                                ]
                            ),
                            html.Br(),
                            dbc.Button(
                                "Filtra Imóveis",
                                color="primary",
                                id="filter_imoveis",
                            ),
                        ]
                    ),
                ]
            ),
            id="fade_filter_imoveis",
            is_in=False,
        ),
    ],
    id="sidebar",
)


def init_app(app: Dash, settings: dict) -> Dash:
    @app.callback(
        Output("sidebar", "className"),
        Output("page-content", "className"),
        Output("side_click", "data"),
        Input("btn_sidebar", "n_clicks"),
        State("side_click", "data"),
    )
    def toggle_sidebar(n, nclick):
        if n:
            if nclick == "SHOW":
                sidebar_style = "hide"
                content_style = ""
                cur_nclick = "HIDDEN"
            else:
                sidebar_style = ""
                content_style = "content_with_sidebar"
                cur_nclick = "SHOW"
        else:
            sidebar_style = ""
            content_style = "content_with_sidebar"
            cur_nclick = "SHOW"

        return sidebar_style, content_style, cur_nclick

    @app.callback(
        Output("page-content", "children"),
        Input("url", "pathname"),
    )
    def render_page_content(pathname):
        if pathname in ["", "/", "/dash", "/dash/table"]:
            return table.layout
        elif pathname == "/dash/map":
            return map.layout

        return dbc.Container(
            [
                html.H1("404: Not found", className="text-danger"),
                html.Hr(),
                html.P(f"The pathname {pathname} was not recognised..."),
            ]
        )

    @app.callback(
        Output("locations", "data"),
        Output("select_local", "options"),
        Output("select_local", "value"),
        Output("fade_search_imoveis", "is_in"),
        Input("search_local", "n_clicks"),
        Input("type_location", "value"),
        State("local", "value"),
    )
    def search_local(_, type_location, value):
        if not value:
            raise PreventUpdate

        locations = Location.parse(
            f"{settings['url']}/locations/?type={type_location}&location={value}")

        locations = [l.dict() for l in locations]
        values = [{"label": l["label"], "value": l["value"]}
                  for l in locations]

        return locations, values, values[0]["value"], True

    @app.callback(
        Output("data", "data"),
        Output("fade_filter_imoveis", "is_in"),
        Output("preco_min", "value"),
        Output("preco_min", "min"),
        Output("preco_min", "max"),
        Output("preco_max", "value"),
        Output("preco_max", "min"),
        Output("preco_max", "max"),
        Output("quarto_min", "value"),
        Output("quarto_min", "min"),
        Output("quarto_min", "max"),
        Output("quarto_max", "value"),
        Output("quarto_max", "min"),
        Output("quarto_max", "max"),
        Input("search_imoveis", "n_clicks"),
        State("locations", "data"),
        State("select_local", "value"),
        State("business_type", "value"),
        State("listing_type", "value"),
    )
    def search_imoveis(_, locations: list[dict], value, business_value, listing_value):
        if not value:
            raise PreventUpdate

        engine = create_engine(
            f"postgresql://{settings['POSTGRES_USER']}:{settings['POSTGRES_PASSWORD']}"
            f"@{settings['POSTGRES_HOST']}:{settings['POSTGRES_PORT']}/{settings['POSTGRES_DB']}"
        )

        selected_location = [l for l in locations if l["value"] == value][0]

        business_type = [
            o for o in depara_tp_contrato if o["value"] == business_value]
        business_type = business_type[0]["value"]

        listing_type = [
            o for o in depara_tp_listings if o["value"] == listing_value]
        listing_type = listing_type[0]["value"]
        LOG.info("Getting data from %s", selected_location)

        url = (
            "{url}/listings?"
            "business_type={business_type}&listing_type={listing_type}&"
            "&locationId={locationId}&city={city}&neighborhood={neighborhood}"
            "&state={state}&stateAcronym={stateAcronym}&zone={zone}"
        ).format(
            url=settings["url"],
            locationId=selected_location["locationId"],
            city=selected_location["city"],
            neighborhood=selected_location["neighborhood"],
            state=selected_location["state"],
            stateAcronym=selected_location["stateAcronym"],
            zone=selected_location["zone"],
            business_type=business_type,
            listing_type=listing_type,
        )
        LOG.info("Getting data from '%s'", url)
        r = get(url)
        r.raise_for_status()

        query = f"""
        SELECT business_type,
               origin,
               url,
               title,
               usable_area,
               unit_types,
               bedrooms,
               bathrooms,
               suites,
               parking_spaces,

               street,
               street_number,
               lat,
               lon,

               price,
               condo_fee,
               (price + condo_fee)           AS total_fee,
               ROUND(predict_total_price, 2) AS total_fee_predict,

               images,
               amenities
        FROM properties
        WHERE
            business_type = '{business_type}'
            AND listing_type = '{listing_type}'
            AND city = '{selected_location['city']}'
            AND neighborhood = '{selected_location['neighborhood']}'
            AND state = '{selected_location['state']}'
            AND state_acronym = '{selected_location['stateAcronym']}'
            AND zone = '{selected_location['zone']}'
            AND active
        """
        LOG.info("Running query: '%s'", query)
        df = pd.read_sql_query(query, engine)
        if df.empty:
            return df.to_dict("records"), True, *[0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]

        return (
            df.to_dict("records"),
            True,
            *min_max(df.total_fee, 500),
            *min_max(df.bedrooms, 1),
        )

    @app.callback(
        Output("filter_description", "children"),
        Input("preco_min", "value"),
        Input("preco_max", "value"),
        Input("quarto_min", "value"),
        Input("quarto_max", "value"),
    )
    def describe_filter(*args):
        if [arg for arg in args if arg is None]:
            raise PreventUpdate

        return "Selecionado imóveis com valores entre [{},{}] e que possuam a quantidade [{},{}] de quartos".format(
            *args
        )

    @app.callback(
        Output("filtered_data", "data"),
        Input("filter_imoveis", "n_clicks"),
        Input("data", "data"),
        State("preco_min", "value"),
        State("preco_max", "value"),
        State("quarto_min", "value"),
        State("quarto_max", "value"),
    )
    def filter_data(n_clicks, data, preco_min, preco_max, quarto_min, quarto_max):
        if n_clicks is None and data is None:
            raise PreventUpdate

        df = pd.DataFrame(data)

        df = df[
            (df["total_fee"] >= preco_min)
            & (df["total_fee"] <= preco_max)
            & (df["bedrooms"] >= quarto_min)
            & (df["bedrooms"] <= quarto_max)
        ]

        return df.to_dict("records")

    return app


LOG = getLogger(__name__)
