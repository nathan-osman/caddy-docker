CWD = $(shell pwd)
PKG = github.com/nathan-osman/caddy-docker
CMD = caddy-docker

SOURCES = $(shell find -type f -name '*.go' ! -path './cache/*')

all: dist/${CMD}

dist/${CMD}: cache dist ${SOURCES}
	docker run \
	    --rm \
	    -e CGO_ENABLED=0 \
	    -v ${CWD}/cache:/go/src \
	    -v ${CWD}/dist:/go/bin \
	    -v ${CWD}:/go/src/${PKG} \
	    -w /go/src/${PKG} \
	    golang:latest \
	    sh -c 'go get $$(go list ./... | grep -v "/cache/")'

cache:
	@mkdir cache

dist:
	@mkdir dist

clean:
	@rm -rf cache dist

.PHONY: clean
