APP_NAME ?= sisyphus
REGISTRY ?= docker.io/vladaluas
VERSION ?= dev

IMAGE := $(REGISTRY)/$(APP_NAME):$(VERSION)

.PHONY: build image push release clean

build:
	go build -o bin/$(APP_NAME) ./cmd/$(APP_NAME)

image:
	docker build -t $(IMAGE) .

push:
	docker push $(IMAGE)

release: image push
	@echo "Published $(IMAGE)"

clean:
	rm -rf bin
