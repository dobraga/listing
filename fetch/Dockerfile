FROM golang:1.18.1 as build

WORKDIR /app

COPY fetch/go.mod ./
COPY fetch/go.sum ./

COPY fetch/./*.go ./
COPY fetch/database/ ./database
COPY fetch/models/ ./models
COPY fetch/property/ ./property
COPY fetch/station/ ./station
COPY fetch/utils/ ./utils

RUN go mod download && CGO_ENABLED=0 go build -o ./go-listing

# RUN-STAGE
FROM gcr.io/distroless/static-debian11

COPY .env /
COPY settings.toml /
COPY --from=build /app/go-listing /

EXPOSE 5000

CMD [ "/go-listing" ]
