.PHONY: build
build:
	go build -ldflags "-s -w" main.go

install:
	sudo chmod +x main
	sudo mv main /usr/local/bin/godoist
