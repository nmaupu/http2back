all: build

fmt:
	go fmt ./...

build:
	env CGO_ENABLED=0 go build

install:
	env CGO_ENABLED=0 go install

clean:
	go clean -i

test:
	go test -v ./...

release:
	env CGO_ENABLED=0 go build -o http2back_linux-amd64

.PHONY: fmt install clean test all
