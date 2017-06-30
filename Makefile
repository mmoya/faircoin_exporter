REPO := mmoya/faircoin-exporter
TAG := latest

build:
	go build
	docker build --pull -t $(REPO):$(TAG) .

push:
	docker push $(REPO):$(TAG)

.PHONY: build-image push-image
