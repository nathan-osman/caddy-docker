CWD = $(shell pwd)
PKG = github.com/nathan-osman/caddy-docker
CMD = caddy-docker

UID = $(shell id -u)
GID = $(shell id -g)

SOURCES = $(shell find -type f -name '*.go' ! -path './cache/*')
BINDATA = $(shell find server/static server/templates)

all: dist/${CMD}

dist/${CMD}: ${SOURCES} server/ab0x.go | cache dist
	docker run \
	    --rm \
	    -e UID=${UID} \
	    -e GID=${GID} \
	    -v ${CWD}/cache:/go/src \
	    -v ${CWD}:/go/src/${PKG} \
	    -w /go/src/${PKG} \
	    nathanosman/bettergo \
	    go get -d ${PKG}/cmd/${CMD}
	docker run \
	    --rm \
	    -e CGO_ENABLED=0 \
	    -e UID=${UID} \
	    -e GID=${GID} \
	    -v ${CWD}/cache:/go/src \
	    -v ${CWD}/dist:/go/bin \
	    -v ${CWD}:/go/src/${PKG} \
	    -w /go/bin \
	    nathanosman/bettergo \
	    go build ${PKG}/cmd/${CMD}

cache:
	@mkdir cache

dist:
	@mkdir dist

server/ab0x.go: ${BINDATA} | dist/fileb0x
	dist/fileb0x b0x.yaml

dist/fileb0x: | dist
	docker run \
	    --rm \
	    -e UID=${UID} \
	    -e GID=${GID} \
	    -v ${CWD}/cache:/go/src \
	    nathanosman/bettergo \
	    go get -d github.com/UnnoTed/fileb0x
	docker run \
	    --rm \
	    -e CGO_ENABLED=0 \
	    -e UID=${UID} \
	    -e GID=${GID} \
	    -v ${CWD}/cache:/go/src \
	    -v ${CWD}/dist:/go/bin \
	    -w /go/bin \
	    nathanosman/bettergo \
	    go build github.com/UnnoTed/fileb0x

clean:
	@rm -f server/ab0x.go
	@rm -rf cache dist

.PHONY: clean
