.PHONY: build build-cli build-relay build-all clean test docker install

VERSION ?= dev
LDFLAGS := -ldflags="-w -s -X main.version=$(VERSION)"

# Build both CLI and relay for current platform
build: build-cli build-relay

build-cli:
	go build $(LDFLAGS) -o bin/filegate ./cmd/filegate

build-relay:
	go build $(LDFLAGS) -o bin/relay ./cmd/relay

# Cross-compile CLI for all platforms
build-all:
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/filegate-linux-amd64 ./cmd/filegate
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/filegate-linux-arm64 ./cmd/filegate
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/filegate-darwin-amd64 ./cmd/filegate
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/filegate-darwin-arm64 ./cmd/filegate
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/filegate-windows-amd64.exe ./cmd/filegate

# Install to GOPATH/bin (makes it available in PATH if GOPATH/bin is in PATH)
install:
	go install $(LDFLAGS) ./cmd/filegate

# Build Docker image for relay
docker:
	docker build -f Dockerfile.relay -t filegate-relay .

# Run relay locally with Docker Compose
docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/ dist/
	go clean
