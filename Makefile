.PHONY: build all publish-dev clean

build:
	mkdir -p dist/
	go build -o dist/ash

all: dist/ash.darwin.amd64 dist/ash.linux.amd64

dist/ash.darwin.amd64:
	mkdir -p dist/
	GOOS=darwin GOARCH=amd64 go build -o dist/ash.darwin.amd64

dist/ash.linux.amd64:
	mkdir -p dist/
	GOOS=linux GOARCH=amd64 go build -o dist/ash.linux.amd64

publish-dev: all
	scripts/unpublish-dev.sh
	scripts/publish-dev.sh dist/ash.darwin.amd64
	scripts/publish-dev.sh dist/ash.linux.amd64

clean:
	rm -rf dist/
