## Build
FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /web-api

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /web-api /web-api

COPY .env ./

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/web-api"]