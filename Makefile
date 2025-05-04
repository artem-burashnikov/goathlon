.PHONY: all test build clean

all: build

build:
	@go build -o goathlon.out ./...

test:
	@go test -v -race ./...

clean:
	@rm goathlon.out
