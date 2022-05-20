# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

ENV GO111MODULE=on

# Install Compile Daemon for go. We'll use it to watch changes in go files
RUN go install github.com/githubnemo/CompileDaemon@latest

ENV GOFLAGS=-mod=vendor

WORKDIR /app

COPY . ./

RUN go mod download
RUN go mod vendor
RUN go mod verify

RUN go build -o /app/src/govwa /app/src/

COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +rx /entrypoint.sh

EXPOSE 8080

WORKDIR /app/src

ENTRYPOINT [ "sh", "/entrypoint.sh" ]
# CMD [ "sh", "/entrypoint.sh" ]
# CMD [ "./govwa" ]