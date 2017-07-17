CWD = $(shell pwd)
PKG = github.com/nathan-osman/caddy-docker
CMD = caddy-docker

SOURCES = $(shell find -type f -name '*.go' ! -path './cache/*')
BINDATA = $(shell find server/static)

all: dist/${CMD}

dist/${CMD}: cache dist ${SOURCES} server/ab0x.go
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

server/ab0x.go: dist/fileb0x
	dist/fileb0x b0x.yaml

dist/fileb0x: dist
	docker run \
	    --rm \
	    -e CGO_ENABLED=0 \
	    -v ${CWD}/cache:/go/src \
	    -v ${CWD}/dist:/go/bin \
	    golang:latest \
	    go get github.com/UnnoTed/fileb0x

clean:
	@rm -f server/ab0x.go
	@rm -rf cache dist

.PHONY: clean
