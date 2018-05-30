BIN=bin
REGISTRY=eu.gcr.io/mirakl-production/kube/admin
HTTP2BACK_VERSION?=0.6
IMAGE_VERSION=$(HTTP2BACK_VERSION)-2

SHORT_NAME=http2back
REMOTE_NAME=$(REGISTRY)/$(SHORT_NAME):$(IMAGE_VERSION)

docker-build:
	docker build --build-arg HTTP2BACK_VERSION=$(HTTP2BACK_VERSION) --no-cache -t mirakl/$(SHORT_NAME):$(IMAGE_VERSION) .

tag: docker-build
	docker tag mirakl/$(SHORT_NAME):$(IMAGE_VERSION) $(REMOTE_NAME)

push: tag
	docker push $(REMOTE_NAME)

all: build

fmt:
	go fmt ./...

build: bin
	env CGO_ENABLED=0 go build -o $(BIN)/http2back

install:
	env CGO_ENABLED=0 go install

clean:
	go clean -i
	rm -rf $(BIN)

test:
	go test -v ./...

bin:
	mkdir -p $(BIN)

.PHONY: fmt install clean test all release
