# Auto generated binary variables helper managed by https://github.com/bwplotka/bingo v0.2.2. DO NOT EDIT.
# All tools are designed to be build inside $GOBIN.
GOPATH ?= $(shell go env GOPATH)
GOBIN  ?= $(firstword $(subst :, ,${GOPATH}))/bin
GO     ?= $(shell which go)

# Bellow generated variables ensure that every time a tool under each variable is invoked, the correct version
# will be used; reinstalling only if needed.
# For example for calens variable:
#
# In your main Makefile (for non array binaries):
#
#include .bingo/Variables.mk # Assuming -dir was set to .bingo .
#
#command: $(CALENS)
#	@echo "Running calens"
#	@$(CALENS) <flags/args..>
#
CALENS := $(GOBIN)/calens-v0.2.0
$(CALENS): .bingo/calens.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/calens-v0.2.0"
	@cd .bingo && $(GO) build -modfile=calens.mod -o=$(GOBIN)/calens-v0.2.0 "github.com/restic/calens"

GOLINT := $(GOBIN)/golint-v0.0.0-20201208152925-83fdc39ff7b5
$(GOLINT): .bingo/golint.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/golint-v0.0.0-20201208152925-83fdc39ff7b5"
	@cd .bingo && $(GO) build -modfile=golint.mod -o=$(GOBIN)/golint-v0.0.0-20201208152925-83fdc39ff7b5 "golang.org/x/lint/golint"

REFLEX := $(GOBIN)/reflex-v0.3.0
$(REFLEX): .bingo/reflex.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/reflex-v0.3.0"
	@cd .bingo && $(GO) build -modfile=reflex.mod -o=$(GOBIN)/reflex-v0.3.0 "github.com/cespare/reflex"

STATICCHECK := $(GOBIN)/staticcheck-v0.1.1
$(STATICCHECK): .bingo/staticcheck.mod
	@# Install binary/ries using Go 1.14+ build command. This is using bwplotka/bingo-controlled, separate go module with pinned dependencies.
	@echo "(re)installing $(GOBIN)/staticcheck-v0.1.1"
	@cd .bingo && $(GO) build -modfile=staticcheck.mod -o=$(GOBIN)/staticcheck-v0.1.1 "honnef.co/go/tools/cmd/staticcheck"

