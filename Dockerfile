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

ENV DB_URL=postgresql://user-manager:qwert12345@db:5432/user-store
ENV SRV_PORT_USER=50051
ENV SRV_NETWORK=tcp
ENV JWT_SECRET=StatusSeeOther

RUN apk update && \
    apk add postgresql-client

RUN apk add --no-cache ca-certificates

WORKDIR /usr/src/app

COPY --from=builder /usr/src/build/user /usr/src/app/user

COPY script/start.sh /start.sh
RUN chmod +x /start.sh

EXPOSE ${SRV_PORT}
