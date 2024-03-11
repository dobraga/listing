ARG BACKEND_PORT

FROM golang:1.21.4 as build

WORKDIR /app

COPY fetch/go.mod ./
COPY fetch/go.sum ./

COPY fetch/cmd/api/main.go ./
COPY fetch/pkg/ ./pkg

RUN go mod download && CGO_ENABLED=0 go build -o ./fetch

# RUN-STAGE
FROM alpine:3.19.1

RUN apk add --update curl \
    && rm -rf /var/cache/apk/*

COPY .env /
COPY settings.toml /
COPY --from=build /app/fetch /bin/fetch

EXPOSE ${BACKEND_PORT}

CMD [ "/bin/fetch" ]
