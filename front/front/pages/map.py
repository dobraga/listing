import pandas as pd
from dash import html, Dash
from folium import Map, Marker
from folium.plugins import MarkerCluster
from dash.exceptions import PreventUpdate
from dash.dependencies import Input, Output

from front.components.box import box

N = 500


layout = html.Div(
    [
        html.Div(id="qtd_sem_latlon"),
        html.Br(),
        html.Iframe(
            id="map",
            srcDoc=Map(height="89%")._repr_html_(),
            width="100%",
            height="800px",
        ),
    ]
)


def init_app(app: Dash) -> Dash:
    @app.callback(
        Output(component_id="qtd_sem_latlon", component_property="children"),
        Output("map", "srcDoc"),
        Input("filtered_data", "data"),
    )
    def create_map(data):
        if not data:
            raise PreventUpdate

        df = pd.DataFrame(data).query("lat != 0 or lon != 0")
        total = df.shape[0]
        if total > N:
            df = df.head(N)
            msg = f"Exibindo {N} imoveis do total de {total} imóveis"
        else:
            msg = f"Exibindo {total} imóveis"

        map = Map(
            location=df[["lat", "lon"]].mean().values,
            height="89%",
            tiles="https://{s}.tile.openstreetmap.de/tiles/osmde/{z}/{x}/{y}.png",
            attr="toner-bcg",
        )

        marker_cluster = MarkerCluster().add_to(map)

        for _, row in df.iterrows():
            Marker(
                location=[row["lat"], row["lon"]],
                popup=box(**row),
            ).add_to(marker_cluster)

        return msg, map._repr_html_()

    return app
