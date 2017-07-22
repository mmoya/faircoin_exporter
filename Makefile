REPO := mmoya/faircoin-exporter
TAG := latest

bin:
	go build

image: bin
	docker build --pull -t $(REPO):$(TAG) .

push:
	docker push $(REPO):$(TAG)

.PHONY: bin image push
