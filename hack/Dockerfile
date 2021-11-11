FROM golang:1.17.2-alpine3.14 as builder
WORKDIR /ch-build/
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY main.go .
COPY cmd /ch-build/cmd
COPY version /ch-build/version

# test build
RUN mkdir -p /ch-build/build && go build -v -o /ch-build/build

# run tests
COPY Makefile .
RUN apk update && apk add make && make test

