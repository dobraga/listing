FROM python:3.10.4-slim-buster

ARG FRONT_PORT

RUN apt-get update && \
    apt-get install -y locales-all

WORKDIR /app

COPY .env ./
COPY front/ ./front
COPY requirements.txt ./

RUN pip3 install -r requirements.txt

CMD gunicorn --bind 0.0.0.0:$FRONT_PORT "front:create_server()" --timeout 0 --log-level error
