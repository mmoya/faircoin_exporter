build-image:
	go build
	docker build --pull -t mmoya/faircoin2-exporter .

.PHONY: build
