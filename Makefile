all: test

test:
	go test -v ./...

build:
	go build

install: build
	go install

server:
	go run main.go server

docker:
	docker build -t scry .

docker-run:
	docker run -d --rm --name scry -p 10080:10080 scry

.PHONY: server test
