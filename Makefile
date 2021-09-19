.PHONY: clean build test docker deploy destroy

## overridable Makefile variables
# test to run
TESTSET = .
# benchmarks to run
BENCHSET ?= .

# version (defaults to short git hash)
VERSION ?= $(shell git rev-parse --short HEAD)

# use correct sed for platform
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
    SED := gsed
else
    SED := sed
endif

PKG_NAME=github.com/Octops/octops-image-syncer
APP_BIN ?= bin/octops-image-syncer
DOCKER_IMAGE_TAG ?= octops/octops-image-syncer:${VERSION}

LDFLAGS := -X "${PKG_NAME}/internal/version.Version=${VERSION}"
LDFLAGS += -X "${PKG_NAME}/internal/version.BuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "${PKG_NAME}/internal/version.GitCommit=$(shell git rev-parse HEAD)"
LDFLAGS += -X "${PKG_NAME}/internal/version.GitBranch=$(shell git rev-parse --abbrev-ref HEAD)"

GO       := GO111MODULE=on GOPRIVATE=github.com/Octops GOSUMDB=off go
GOBUILD  := CGO_ENABLED=0 $(GO) build $(BUILD_FLAG)
GOTEST   := $(GO) test -gcflags='-l' -p 3

CURRENT_DIR := $(shell pwd)
FILES    := $(shell find internal cmd -name '*.go' -type f -not -name '*.pb.go' -not -name '*_generated.go' -not -name '*_test.go')
TESTS    := $(shell find internal cmd -name '*.go' -type f -not -name '*.pb.go' -not -name '*_generated.go' -name '*_test.go')

default: clean build

clean:
	rm -rf bin/*

$(APP_BIN):
	go build -ldflags '$(LDFLAGS)' -a -installsuffix cgo -o $(APP_BIN) .

build: clean $(APP_BIN)

test:
	$(GOTEST) -run=$(TESTSET) ./...
	@echo
	@echo Configured tests ran ok.

test-strict:
	$(GO) test -p 3 -run=$(TESTSET) -gcflags='-l -m' -race ./...
	@echo
	@echo Configured tests ran ok.

vendor:
	$(GO) mod vendor

docker:
	docker build -t $(DOCKER_IMAGE_TAG) .

push: docker
	docker push $(DOCKER_IMAGE_TAG)

deploy: docker
	docker save --output bin/octops-image-syncer-v0.0.1.tar $(DOCKER_IMAGE_TAG)
	rsync -v bin/octops-image-syncer-v0.0.1.tar ${SSH_REMOTE}:/home/octops/
	ssh ${SSH_REMOTE} sudo -S k3s ctr images import /home/octops/octops-image-syncer-v0.0.1.tar
	kubectl apply -f hack/install.yaml

destroy:
	kubectl delete -f hack/install.yaml