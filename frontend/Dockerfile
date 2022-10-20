## Build
FROM golang:1.17-buster AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /frontend

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /frontend /frontend
COPY templates /templates

EXPOSE 8090

USER nonroot:nonroot

ENTRYPOINT ["/frontend"]    