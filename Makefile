.PHONY: build run install clean test

build:
	go build -o build/droid-config ./cmd/droid-config

run: build
	./build/droid-config

install: build
	cp build/droid-config ~/.local/bin/droid-config

clean:
	rm -rf build/

test:
	go test ./...
	go vet ./...
