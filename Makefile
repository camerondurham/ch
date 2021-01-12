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
# removing 386 as a target architecture
ARCHITECTURES=amd64

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

default: build

all: clean build_all install

build:
	go build ${LDFLAGS} -o ${BINARY}

build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell mkdir -p dist/$(BINARY)-$(GOOS)-$(GOARCH); GOOS=$(GOOS) GOARCH=$(GOARCH) go build -v -o dist/$(BINARY)-$(GOOS)-$(GOARCH))))

TO_ZIP_DIRS = $(filter %/, $(wildcard dist/*/))  	# Find all directories in static/projects
TO_ZIP_NAMES = $(patsubst %/,%,$(TO_ZIP_DIRS))  	# Remove trailing /
ZIP_TARGETS = $(addsuffix .zip,$(TO_ZIP_NAMES))  	# Add .zip

debug: build_all
	@echo $(TO_ZIP_DIRS)
	@echo $(TO_ZIP_NAMES)
	@echo $(ZIP_TARGETS)

$(ZIP_TARGETS):
	cd $(basename $@)/.. && zip -FSr $(notdir $@) $(notdir $(basename $@))

# edit .github/workflows/build-and-release.yml if this name changes
zip_exe: $(ZIP_TARGETS)

install:
	go install ${LDFLAGS}

# Remove only what we've created
clean:
	find ${ROOT_DIR} -name '${BINARY}[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' | xargs rm -rf

.PHONY: check clean install build_all all zip_exe

