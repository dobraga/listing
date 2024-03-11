ARG MODEL_PORT

FROM python:3.10.13-slim

RUN apt update -y && apt install -y libgomp1 curl

WORKDIR /app

COPY .env ./
COPY settings.toml ./
COPY model/api/ ./api
COPY model/model/ ./model
COPY model/requirements.txt ./

RUN pip install --no-cache-dir -r requirements.txt

CMD uvicorn --host 0.0.0.0 --port $MODEL_PORT "api:app" --log-level warning
