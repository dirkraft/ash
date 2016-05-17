.PHONY: clean build build-all

build:
	mkdir -p dist/
	go build -o dist/ash

all:
	mkdir -p dist/
	mkdir -p dist/
	GOOS=darwin GOARCH=amd64 go build -o dist/ash.darwin.amd64
	GOOS=linux GOARCH=amd64 go build -o dist/ash.linux.amd64

clean:
	rm -rf dist/
