export GOPATH := $(shell pwd)
default: build

init:
	@rm -f bin/main
	@cd src/main && go get

build: init
	@go build -o bin/main src/main/main.go 

run: build
	@bin/main
