import pandas as pd
from dash import Dash
from math import ceil
from dash.dash_table import DataTable
from dash.dependencies import Input, Output


layout = DataTable(
    id="table",
    columns=[
        {"name": "Título", "id": "title", "type": "text", "presentation": "markdown"},
        {"name": "Área Útil", "id": "usable_area"},
        {"name": "Tipo", "id": "unit_types"},
        # {"name": "Estação", "id": "estacao"},
        # {"name": "Distância", "id": "distance"},
        {"name": "Preço Total", "id": "total_fee"},
        {"name": "Preço Predito", "id": "total_fee_predict"},
        {"name": "Erro", "id": "error"},
        {"name": "Possui Loc", "id": "fl_latlon"},
    ],
    style_header={"fontWeight": "bold"},
    page_current=0,
    page_size=25,
    page_action="custom",
    sort_action="custom",
    sort_mode="multi",
    sort_by=[],
    style_cell={"textOverflow": "ellipsis"},
    style_cell_conditional=[
        {"if": {"column_id": "title"}, "maxWidth": 45, "textAlign": "left"},
        {"if": {"column_id": "usable_area"}, "maxWidth": 15},
        {"if": {"column_id": "unit_types"}, "maxWidth": 20, "textAlign": "left"},
        {"if": {"column_id": "distance"}, "maxWidth": 20},
        {"if": {"column_id": "total_fee"}, "maxWidth": 20},
        {"if": {"column_id": "total_fee_predict"}, "maxWidth": 20},
        {"if": {"column_id": "error"}, "maxWidth": 20},
        {"if": {"column_id": "fl_latlon"}, "maxWidth": 20},
    ],
)


def init_app(app: Dash) -> Dash:
    cols = ["title", "usable_area", "unit_types"]
    cols += ["total_fee", "total_fee_predict", "error", "fl_latlon"]

    @app.callback(
        Output("table", "data"),
        Output("table", "page_count"),
        Input("filtered_data", "data"),
        Input("table", "sort_by"),
        Input("table", "page_current"),
        Input("table", "page_size"),
    )
    def updateTable(data, sort_by, page_current, page_size):
        dff = pd.DataFrame(data)

        if dff.shape[0] == 0:
            return [], 1

        dff["error"] = (dff["total_fee"] - dff["total_fee_predict"]).round(0)
        dff["total_fee_predict"] = dff["total_fee_predict"].round(0)
        dff["fl_latlon"] = ((dff["lat"] != 0) & (dff["lon"] != 0)).astype(int)

        dff = dff.sort_values("fl_latlon", ascending=False)\
            .groupby(["title", "total_fee"], sort=False)\
            .head(1)

        dff["title"] = dff[["title", "origin", "url"]].apply(
            lambda x: f"[{x[0]}](https://www.{x[1]}.com.br/imovel/{x[2]})", axis=1
        )

        if len(sort_by):
            dff = dff.sort_values(
                [col["column_id"] for col in sort_by],
                ascending=[col["direction"] == "asc" for col in sort_by],
            )
        else:
            dff = dff.sort_values("error")

        dff = dff[cols]

        return (
            dff.iloc[page_current * page_size: (page_current + 1) * page_size].to_dict(
                "records"
            ),
            ceil(dff.shape[0] / page_size),
        )

    return app
