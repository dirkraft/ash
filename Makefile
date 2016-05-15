.PHONY: clean build build-all

build:
	mkdir -p dist/
	go build -o dist/ash

clean:
	rm -rf dist/
