FROM golang:1.17.2-alpine3.14 as builder
WORKDIR /ch-build/
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY main.go .
COPY cmd /ch-build/cmd
COPY version /ch-build/version
RUN mkdir -p /ch-build/build && go build -v -o /ch-build/build

FROM alpine:3.14.2 as cli
WORKDIR /cli/
COPY --from=builder /ch-build/build ./
ENTRYPOINT ["/cli/ch"]
