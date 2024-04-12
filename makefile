.PHONY: run build test

build:
	go build

run: build
	./seed-go

test:
	go test -v ./...
