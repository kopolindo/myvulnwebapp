# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor

WORKDIR /app

COPY . ./

RUN go mod download
RUN go mod vendor
RUN go mod verify

RUN go build -o /app/src/govwa /app/src/

EXPOSE 8080

WORKDIR /app/src

CMD [ "./govwa" ]