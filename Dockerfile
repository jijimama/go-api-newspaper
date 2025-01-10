FROM golang:1.23.2-alpine3.19

RUN apk --no-cache add curl

WORKDIR /go/src/web

COPY .. .

RUN go mod download

ENV GO111MODULE=on

RUN go build main.go

ENTRYPOINT ./main
