include golang.mk
.DEFAULT_GOAL := test # override default goal set in library makefile

.PHONY: all test build vendor
SHELL := /bin/bash
PKG := github.com/Clever/oplog-dump/cmd/oplog-dump
PKGS := $(shell go list ./... | grep -v /vendor)
EXECUTABLE := $(shell basename $(PKG))
$(eval $(call golang-version-check,1.9))

all: test build

build:
	go build -o bin/$(EXECUTABLE) $(PKG)

test: $(PKGS)
$(PKGS): golang-test-all-deps
	$(call golang-test-all,$@)



install_deps: golang-dep-vendor-deps
	$(call golang-dep-vendor)