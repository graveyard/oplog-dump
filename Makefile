include golang.mk
.DEFAULT_GOAL := test # override default goal set in library makefile

.PHONY: all test build vendor
SHELL := /bin/bash
PKG := github.com/Clever/oplog-dump/cmd/oplog-dump
PKGS := $(shell go list ./... | grep -v /vendor)
EXECUTABLE := $(shell basename $(PKG))
$(eval $(call golang-version-check,1.8))

all: test build

build:
	go build -o bin/$(EXECUTABLE) $(PKG)

test: $(PKGS)
$(PKGS): golang-test-all-deps
	$(call golang-test-all,$@)

vendor: golang-godep-vendor-deps
	$(call golang-godep-vendor,$(PKGS))
