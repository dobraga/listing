import pandas as pd
from dash import html, Dash
from folium import Map, Marker
from folium.plugins import MarkerCluster
from dash.exceptions import PreventUpdate
from dash.dependencies import Input, Output

from front.components.box import box


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

        df = pd.DataFrame(data)
        total = df.shape[0]
        qtd_sem_latlon = (df[["lat", "lon"]].max(axis=1) != 0).sum()
        df = df.query("lat != 0 or lon != 0").head(150)

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

        return (
            f"Exibindo 150 imoveis do total de {total} imóveis sendo {qtd_sem_latlon} sem latitude ou longitude",
            map._repr_html_(),
        )

    return app
