CWD = $(shell pwd)
PKG = github.com/nathan-osman/caddy-docker
CMD = caddy-docker

SOURCES = $(shell find -type f -name '*.go')

all: dist/${CMD}

dist/${CMD}: dist
	docker run \
	    --rm \
	    -e CGO_ENABLED=0 \
	    -v ${CWD}:/go/src/${PKG} \
	    -v ${CWD}/dist:/go/bin \
	    -w /go/src/${PKG} \
	    golang:latest \
	    go get ./...

dist:
	@mkdir dist

clean:
	@rm -rf dist

.PHONY: clean
