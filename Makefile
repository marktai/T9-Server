export GOPATH := $(shell pwd)
default: run

init:
	rm bin/main
	cd src/main && go get

run: init
	go build -o bin/main src/main/main.go 
	bin/main
