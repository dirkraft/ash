.PHONY: build all publish-dev clean

build:
	mkdir -p dist/
	go build -o dist/ash ash.go
	go build -o dist/acp acp.go

all: dist/ash.darwin.amd64 dist/ash.linux.amd64

dist/ash.darwin.amd64:
	mkdir -p dist/
	GOOS=darwin GOARCH=amd64 go build -o dist/ash.darwin.amd64 ash.go
	GOOS=darwin GOARCH=amd64 go build -o dist/acp.darwin.amd64 acp.go

dist/ash.linux.amd64:
	mkdir -p dist/
	GOOS=linux GOARCH=amd64 go build -o dist/ash.linux.amd64 ash.go
	GOOS=linux GOARCH=amd64 go build -o dist/acp.linux.amd64 acp.go

publish-dev: all
	scripts/unpublish-dev.sh
	scripts/publish-dev.sh dist/ash.darwin.amd64
	#scripts/publish-dev.sh dist/acp.darwin.amd64
	scripts/publish-dev.sh dist/ash.linux.amd64
	#scripts/publish-dev.sh dist/acp.linux.amd64

clean:
	rm -rf dist/
