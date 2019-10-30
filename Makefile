.PHONY: build push
TAG?=latest

build:
	operator-sdk build sinoreps/cmpp-operator:$(TAG)

push:
	docker push sinoreps/cmpp-operator:$(TAG)
