FROM golang:alpine

WORKDIR /app

COPY .env ./
COPY settings.toml ./
COPY fetch/go.mod ./
COPY fetch/go.sum ./
RUN go mod download

COPY fetch/*.go ./

RUN go build -o /go-listing

EXPOSE 5000

CMD [ "/go-listing" ]
