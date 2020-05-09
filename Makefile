# gcoreclient Makefile

PWD := $(shell pwd)
BASE_DIR := $(shell basename $(PWD))
# Keep an existing GOPATH, make a private one if it is undefined
GOPATH_DEFAULT := $(PWD)/.go
export GOPATH ?= $(GOPATH_DEFAULT)
GOBIN_DEFAULT := $(GOPATH)/bin
export GOBIN ?= $(GOBIN_DEFAULT)
export GO111MODULE := on
TESTARGS_DEFAULT := -v -race
TESTARGS ?= $(TESTARGS_DEFAULT)
PKG := $(shell awk '/^module/ { print $$2 }' go.mod)
DEST := $(GOPATH)/src/$(GIT_HOST)/$(BASE_DIR)
SOURCES := $(shell find $(DEST) -name '*.go' 2>/dev/null)
HAS_GOLANGCI := $(shell command -v golangci-lint;)
HAS_GOIMPORTS := $(shell command -v goimports;)

TARGETS		?= darwin/amd64 linux/amd64 linux/386 linux/arm linux/arm64 linux/ppc64le linux/s390x
DIST_DIRS	= find * -type d -exec

TEMP_DIR	:=$(shell mktemp -d)

GOOS		?= $(shell go env GOOS)
VERSION		?= $(shell git describe --exact-match 2> /dev/null || \
			   git describe --match=$(git rev-parse --short=8 HEAD) --always --dirty --abbrev=8)
GOARCH		:= amd64
TAGS		:=
LDFLAGS		:= "-w -s"
CMD_PACKAGE := ./gcoreclient
BINARY 		:= ./gcoreclient/gcoreclient

# CTI targets

$(GOBIN):
	echo "create gobin"
	mkdir -p $(GOBIN)

work: $(GOBIN)

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
	-ldflags $(LDFLAGS) \
	-o $(BINARY) \
	$(CMD_PACKAGE)

install: check test
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go install \
	-ldflags $(LDFLAGS) \
	$(CMD_PACKAGE)

test: unit functional

check: work fmt vet goimports golangci
unit: work
	go test -tags=unit $(TESTARGS) ./...

functional:
	@echo "$@ not yet implemented"

fmt:
	go fmt ./...

goimports:
ifndef HAS_GOIMPORTS
	echo "installing goimports"
	GO111MODULE=off go get golang.org/x/tools/cmd/goimports
endif
	goimports -d $(shell find . -iname "*.go")

vet:
	go vet ./...

golangci:
ifndef HAS_GOLANGCI
	GO111MODULE=off go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.26.0
endif
	golangci-lint run ./...

cover: work
	go test -tags=unit ./... -cover


prepare:
ifndef HAS_GOLANGCI
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.26.0
endif
	echo "golangci-lint already installed"
ifndef HAS_GOIMPORTS
	echo "installing goimports"
	GO111MODULE=off go get golang.org/x/tools/cmd/goimports
endif
	echo "goimports already installed"

env:
	@echo "PWD: $(PWD)"
	@echo "BASE_DIR: $(BASE_DIR)"
	@echo "GOPATH: $(GOPATH)"
	@echo "GOROOT: $(GOROOT)"
	@echo "DEST: $(DEST)"
	@echo "PKG: $(PKG)"
	go version
	go env

shell:
	$(SHELL) -i

clean: work
	rm -rf $(BINARY)

version:
	@echo ${VERSION}

.PHONY: bindep install build cover work fmt functional lint test version clean prepare
