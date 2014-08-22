all: build test

deps:
	go get -t ./...

build: deps
	go build -o bin/oplog-dump github.com/Clever/oplog-dump/cmd/oplog-dump

test: deps
	go test github.com/Clever/oplog-dump/cmd/oplog-dump
