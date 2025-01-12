.PHONY: build clean test

build:
	go build -o piecehub ./cmd/piecehub/...

clean:
	rm -rf piecehub

test:
	go test -v ./...
