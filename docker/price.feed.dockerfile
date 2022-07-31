# syntax=docker/dockerfile:1
FROM golang:1.18-bullseye AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod verify

COPY app ./app
COPY cmd ./cmd
COPY pkg ./pkg

RUN go build -o /build/hmq.prices.contract_feed ./cmd/contract_feed/main.go

FROM debian:bullseye-slim

WORKDIR /

RUN apt update && apt -y upgrade && apt -y install ca-certificates && apt -y autoremove

COPY --from=build /build/hmq.prices.contract_feed /usr/local/bin/hmq.prices.contract_feed
COPY spec /usr/local/share

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/hmq.prices.contract_feed"]

