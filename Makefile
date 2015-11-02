SHELL := /bin/bash
GOVERSION := $(shell go version | grep 1.5)
ifeq "$(GOVERSION)" ""
  $(error must be running Go version 1.5)
endif

export GO15VENDOREXPERIMENT = 1

all: build test

deps:
	go get -t ./...

build: deps
	go build -o bin/oplog-dump github.com/Clever/oplog-dump/cmd/oplog-dump

test: deps
	go test github.com/Clever/oplog-dump/cmd/oplog-dump
