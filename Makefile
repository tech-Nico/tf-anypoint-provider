VERSION=1.0.0
BUILDBOX_TAG ?= golang:1.9.0-stretch

default: build

build: 
	docker build . -t tf-anypoint:$(VERSION)

.PHONY: build 