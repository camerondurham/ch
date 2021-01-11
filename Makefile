# ch: container-helper
# using makefile template from: https://gist.github.com/cjbarker/5ce66fcca74a1928a155cfb3fea8fac4

# Check for required command tools to build or stop immediately
EXECUTABLES = git go find pwd
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

BINARY=ch
VERSION=0.0.1
BUILD=`git rev-parse HEAD`
PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

default: build

all: clean build_all install

build:
	go build ${LDFLAGS} -o ${BINARY}

build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -v -o $(BINARY)-$(GOOS)-$(GOARCH))))

install:
	go install ${LDFLAGS}

# Remove only what we've created
clean:
	find ${ROOT_DIR} -name '${BINARY}[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' -delete

.PHONY: check clean install build_all all


# all: $(DISTS)

# BUILDS=\
#   darwin-amd64  \
#   linux-386     \
#   linux-amd64   \
#   linux-arm     \
#   linux-arm64   \
#   windows-386   \
#   windows-amd64 \

# dist:
# 	@mkdir -p dist

# $(DISTS): OS = $(word 1,$(subst -, ,$*))
# $(DISTS): ARCH = $(word 2,$(subst -, ,$*))
# $(DISTS): DIST = "$(OS)-$(ARCH)"

# $(DISTS): dist/$(NAME)-%-$(VERSION).tgz: dist
# 	@echo "building: $@"
# 	@echo "OS = $(OS)"
# 	@echo "ARCH = $(ARCH)"
# 	@mkdir -p "dist/$(DIST)"
# 	@GOOS=$(OS) GOARCH=$(ARCH) go build -o "dist/$(DIST)"

# clean:
# 	rm -rf dist

# cleanup:
# 	go mod tidy

# get-tools:
# 	go get -t \
# 	github.com/spf13/cobra/cobra

# todo:
# 	 git grep -EI "TODO|FIXME"
