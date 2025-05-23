FROM golang:1.24.1 AS builder

LABEL stage=builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /usr/src/build

ADD go.mod ./
ADD go.sum ./
RUN go mod download

COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./cmd ./cmd

RUN go build -o user ./cmd/app/main.go

FROM alpine:latest

LABEL authors="ekvo"

ENV DB_HOST=db
ENV DB_PORT=5432
ENV DB_USER=manager
ENV DB_PASSWORD=qwert12345
ENV DB_PORT=5432
ENV DB_NAME=userdb
ENV DB_MAX_CONN=10
ENV DB_MIN_CONN=1
ENV DB_CONN_MAX_LIFE_TIME=24h
ENV DB_CONN_MAX_IDLE_TIME=15m
ENV DB_CONN_TIMEOUT=1m
ENV DB_HEALTH_CHECK_PERIOD=1m

ENV MIGRATION_PATH=sql/migrations

ENV SRV_PORT=50051
ENV SRV_NETWORK=tcp

ENV JWT_SECRET=StatusSeeOther

RUN apk update && \
    apk add postgresql-client

RUN apk add --no-cache ca-certificates

WORKDIR /usr/src/app

COPY --from=builder /usr/src/build/user /usr/src/app/user
COPY ./sql /usr/src/app/sql

COPY script/start.sh /start.sh
RUN chmod +x /start.sh

EXPOSE ${SRV_PORT}
