.PHONY: all build ci-bot clean

all:build

build:ci-bot

build-image:ci-bot-image

ci-bot:
	go build -o ci-bot ./cmd/cibot

ci-bot-image:
	docker build -t openeuler/cibot:latest ./	

clean:
	rm -rf ./ci-bot
