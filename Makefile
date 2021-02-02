include .bingo/Variables.mk

SHELL := bash
NAME := vcloud-csi-driver
IMPORT := github.com/proact-de/$(NAME)
BIN := bin

EXECUTABLE := $(NAME)
UNAME := $(shell uname -s)

ifeq ($(UNAME), Darwin)
	GOBUILD ?= CGO_ENABLED=0 go build -i
else
	GOBUILD ?= CGO_ENABLED=0 go build
endif

PACKAGES ?= $(shell go list ./...)
SOURCES ?= $(shell find . -name "*.go" -type f)
GENERATE ?= $(PACKAGES)

TAGS ?= netgo

ifndef OUTPUT
	ifneq ($(DRONE_TAG),)
		OUTPUT ?= $(subst v,,$(DRONE_TAG))
	else
		OUTPUT ?= testing
	endif
endif

ifndef VERSION
	ifneq ($(DRONE_TAG),)
		VERSION ?= $(subst v,,$(DRONE_TAG))
	else
		VERSION ?= $(shell git rev-parse --short HEAD)
	endif
endif

ifndef DATE
	DATE := $(shell date -u '+%Y%m%d')
endif

ifndef SHA
	SHA := $(shell git rev-parse --short HEAD)
endif

LDFLAGS += -s -w -extldflags "-static" -X "$(IMPORT)/pkg/version.String=$(VERSION)" -X "$(IMPORT)/pkg/version.Revision=$(SHA)" -X "$(IMPORT)/pkg/version.Date=$(DATE)"
GCFLAGS += all=-N -l

.PHONY: all
all: build

.PHONY: sync
sync:
	go mod download

.PHONY: clean
clean:
	go clean -i ./...
	rm -rf $(BIN)

.PHONY: fmt
fmt:
	gofmt -s -w $(SOURCES)

.PHONY: vet
vet:
	go vet $(PACKAGES)

.PHONY: staticcheck
staticcheck: $(STATICCHECK)
	$(STATICCHECK) -tags '$(TAGS)' $(PACKAGES)

.PHONY: lint
lint: $(GOLINT)
	for PKG in $(PACKAGES); do $(GOLINT) -set_exit_status $$PKG || exit 1; done;

.PHONY: generate
generate:
	go generate $(GENERATE)

.PHONY: changelog
changelog: $(CALENS)
	$(CALENS) >| CHANGELOG.md

.PHONY: test
test: $(GOVERAGE)
	$(GOVERAGE) -v -coverprofile coverage.out $(PACKAGES)

.PHONY: build
build: $(BIN)/$(EXECUTABLE)

$(BIN)/$(EXECUTABLE): $(SOURCES)
	$(GOBUILD) -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $@ ./cmd/$(NAME)

$(BIN)/$(EXECUTABLE)-debug: $(SOURCES)
	$(GOBUILD) -v -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -gcflags '$(GCFLAGS)' -o $@ ./cmd/$(NAME)
