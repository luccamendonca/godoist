.PHONY: build
build:
	go build -ldflags "-s -w" main.go
