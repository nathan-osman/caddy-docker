CWD = $(shell pwd)
PKG = github.com/nathan-osman/caddy-docker
CMD = caddy-docker

UID = $(shell id -u)
GID = $(shell id -g)

SOURCES = $(shell find -type f -name '*.go' ! -path './cache/*')
BINDATA = $(shell find server/static)

all: dist/${CMD}

dist/${CMD}: ${SOURCES} server/ab0x.go | cache dist
	docker run \
	    --rm \
	    -e CGO_ENABLED=0 \
	    -v ${CWD}/cache:/go/src \
	    -v ${CWD}/dist:/go/bin \
	    -v ${CWD}:/go/src/${PKG} \
	    -w /go/src/${PKG} \
	    -e UID=${UID} \
	    -e GID=${GID} \
	    nathanosman/bettergo \
	    sh -c 'go get $$(go list ./... | grep -v "/cache/")'

cache:
	@mkdir cache

dist:
	@mkdir dist

server/ab0x.go: ${BINDATA} | dist/fileb0x
	dist/fileb0x b0x.yaml

dist/fileb0x: | dist
	docker run \
	    --rm \
	    -e CGO_ENABLED=0 \
	    -v ${CWD}/cache:/go/src \
	    -v ${CWD}/dist:/go/bin \
	    -e UID=${UID} \
	    -e GID=${GID} \
	    nathanosman/bettergo \
	    go get github.com/UnnoTed/fileb0x

clean:
	@rm -f server/ab0x.go
	@rm -rf cache dist

.PHONY: clean
