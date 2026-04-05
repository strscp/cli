.PHONY: build test lint clean snapshot

build:
	go build -o bin/starscope-cli ./cmd/starscope-cli

test:
	go test -race -count=1 ./...

test-integration:
	go test -race -count=1 -tags=integration ./...

lint:
	golangci-lint run

coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf bin/ dist/

snapshot:
	goreleaser release --snapshot --clean
