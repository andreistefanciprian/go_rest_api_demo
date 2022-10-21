## Build
FROM golang:1.17-buster AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /web-api

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /web-api /web-api

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/web-api"]