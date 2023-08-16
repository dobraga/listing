FROM golang:1.18.1 as build

WORKDIR /app

COPY fetch/go.mod ./
COPY fetch/go.sum ./

COPY fetch/cmd/api/main.go ./
COPY fetch/pkg/ ./pkg

RUN go mod download && CGO_ENABLED=0 go build -o ./go-listing

# RUN-STAGE
FROM gcr.io/distroless/static-debian11

COPY .env /
COPY settings.toml /
COPY --from=build /app/go-listing /

EXPOSE 5000

CMD [ "/go-listing" ]
