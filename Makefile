.PHONY: develop build all publish-dev clean

develop:
	go get ./ash ./acp

build:
	mkdir -p dist/
	go build -o dist/ash ash/main.go
	go build -o dist/acp acp/main.go

all: dist/ash.darwin.amd64 dist/ash.linux.amd64 dist/ash.windows.amd64

dist/ash.darwin.amd64:
	mkdir -p dist/
	GOOS=darwin GOARCH=amd64 go build -o dist/ash.darwin.amd64 ash/main.go
	GOOS=darwin GOARCH=amd64 go build -o dist/acp.darwin.amd64 acp/main.go

dist/ash.linux.amd64:
	mkdir -p dist/
	GOOS=linux GOARCH=amd64 go build -o dist/ash.linux.amd64 ash/main.go
	GOOS=linux GOARCH=amd64 go build -o dist/acp.linux.amd64 acp/main.go

dist/ash.windows.amd64:
	mkdir -p dist/
	GOOS=windows GOARCH=amd64 go build -o dist/ash.windows.amd64 ash/main.go
	GOOS=windows GOARCH=amd64 go build -o dist/acp.windows.amd64 acp/main.go

publish-dev: all
	scripts/unpublish-dev.sh
	scripts/publish-dev.sh dist/ash.darwin.amd64
	#scripts/publish-dev.sh dist/acp.darwin.amd64
	scripts/publish-dev.sh dist/ash.linux.amd64
	#scripts/publish-dev.sh dist/acp.linux.amd64
	scripts/publish-dev.sh dist/ash.windows.amd64
	#scripts/publish-dev.sh dist/acp.windows.amd64
	scripts/update-dev-tag.sh

clean:
	rm -rf dist/
