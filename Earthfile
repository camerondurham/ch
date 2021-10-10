FROM public.ecr.aws/bitnami/golang:1.17
WORKDIR /ch-build

deps:
    COPY go.mod go.sum .
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

build:
    FROM +deps
    COPY main.go .
    COPY --dir cmd version .
    RUN mkdir -p build/ch-build && go build -v -o build/ch-build
    SAVE ARTIFACT build/ch-build /ch-build AS LOCAL build/ch-build

docker:
    COPY +build/ch-build .
    ENTRYPOINT ["/ch-build/ch"]
    SAVE IMAGE ch-docker

all:
    BUILD +build
    BUILD +docker