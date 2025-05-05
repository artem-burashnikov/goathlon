.PHONY: all test build clean

all: build

build:
	@go build -o goathlon.out ./...

test:
	@go test -v -race ./...

coverage:
	@go test -coverprofile=coverage.out -covermode=atomic -race ./...
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out

clean:
	@rm -f goathlon.out coverage.out coverage.html
