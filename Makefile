.PHONY: build run

build:
    @echo "Building Go application..."
    @cd src && go build -o main

run: build
    @echo "Running Go application..."
    @./src/main