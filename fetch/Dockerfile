FROM golang:1.18.1 as build

WORKDIR /app

COPY .env ./
COPY settings.toml ./
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./*.go ./
COPY database/ ./database
COPY models/ ./models
COPY property/ ./property
COPY station/ ./station
COPY utils/ ./utils

RUN CGO_ENABLED=0 go build -o ./go-listing

# RUN-STAGE
FROM gcr.io/distroless/static-debian11

COPY --from=build /app/go-listing /
COPY --from=build /app/.env ./
COPY --from=build /app/settings.toml ./

EXPOSE 5000

CMD [ "/go-listing" ]
