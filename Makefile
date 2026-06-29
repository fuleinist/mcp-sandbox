.PHONY: build test lint clean

BINARY=mcp-sandbox

build:
	go build -o $(BINARY) ./cmd/

test:
	go test -v -race -count=1 ./...

lint:
	golangci-lint run

clean:
	rm -f $(BINARY) $(BINARY).exe
