.PHONY: build

all: build

build:
	@go build -o bin -v ./server/main.go

