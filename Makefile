.PHONY: build push
TAG?=latest

build:
	docker build -t sinoreps/cmpp-operator:$(TAG) . -f build/Dockerfile

push:
	docker push sinoreps/cmpp-operator:$(TAG)
