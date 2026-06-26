BINARY := main

.PHONY: build lint run

build:
	go build -o $(BINARY) .

lint:
	golangci-lint run

run: build
	./$(BINARY)
