SHELL := /bin/bash
PKG := github.com/Clever/oplog-dump/cmd/oplog-dump
PKGS := $(shell go list ./... | grep -v /vendor)
EXECUTABLE := oplog-dump
.PHONY: all test build vendor

GOVERSION := $(shell go version | grep 1.5)
ifeq "$(GOVERSION)" ""
  $(error must be running Go version 1.5)
endif
export GO15VENDOREXPERIMENT = 1

GOLINT := $(GOPATH)/bin/golint
$(GOLINT):
	go get github.com/golang/lint/golint

GODEP := $(GOPATH)/bin/godep
$(GODEP):
	go get -u github.com/tools/godep

all: build test

build:
	go build -o bin/$(EXECUTABLE) $(PKG)

test:
	go test $(PKGS)

vendor: $(GODEP)
	$(GODEP) save $(PKGS)
	find vendor/ -path '*/vendor' -type d | xargs -IX rm -r X # remove any nested vendor directories
