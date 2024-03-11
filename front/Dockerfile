ARG FRONT_PORT

FROM python:3.10.13-slim

RUN apt-get update && \
    apt-get install -y locales-all

WORKDIR /app

COPY .env ./
COPY settings.toml ./
COPY front/front/ ./front
COPY front/requirements.txt ./

RUN pip3 install -r requirements.txt

CMD gunicorn --bind 0.0.0.0:$FRONT_PORT "front:create_server()" --timeout 0 --log-level error
