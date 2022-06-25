FROM golang:alpine
RUN apk add build-base

WORKDIR /app

COPY .env ./
COPY settings.toml ./
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./*.go ./

RUN go build -o /go-listing

EXPOSE 5000

CMD [ "/go-listing" ]
