REPO := mmoya/faircoin-exporter
TAG := latest

bin:
	go build

docker-bin:
	docker run --rm -it \
	    -v "$$GOPATH/bin":/gopath/bin \
	    -v "$$GOPATH/pkg":/gopath/pkg \
	    -v "$$GOPATH/src":/gopath/src \
	    -v "$$PWD":/app \
	    -e "GOPATH=/gopath" \
	    -e "CGO_ENABLED=1" \
	    -w /app \
	    mmoya/golang-zeromq:1.7-alpine3.6 \
	        go build -a --installsuffix cgo --ldflags="-s" -o faircoin_exporter.docker

image: docker-bin
	docker build --pull -t $(REPO):$(TAG) .

push:
	docker push $(REPO):$(TAG)

.PHONY: bin docker-bin image push
