CWD = $(shell pwd)
PKG = github.com/nathan-osman/caddy-docker
CMD = caddy-docker

UID = $(shell id -u)
GID = $(shell id -g)

SOURCES = $(shell find -type f -name '*.go' ! -path './cache/*')
BINDATA = $(shell find server/static server/templates)

all: dist/${CMD}

dist/${CMD}: ${SOURCES} server/ab0x.go | cache dist
	@docker run \
	    --rm \
	    -e CGO_ENABLED=0 \
	    -e UID=${UID} \
	    -e GID=${GID} \
	    -v ${CWD}/cache/lib:/go/lib \
	    -v ${CWD}/cache/src:/go/src \
	    -v ${CWD}/dist:/go/bin \
	    -v ${CWD}:/go/src/${PKG} \
	    nathanosman/bettergo \
	    go get -pkgdir /go/lib ${PKG}/cmd/${CMD}
	@touch dist/${CMD}

server/ab0x.go: ${BINDATA} dist/fileb0x b0x.yaml
	dist/fileb0x b0x.yaml

dist/fileb0x: | cache dist
	@docker run \
	    --rm \
	    -e CGO_ENABLED=0 \
	    -e UID=${UID} \
	    -e GID=${GID} \
	    -v ${CWD}/cache/lib:/go/lib \
	    -v ${CWD}/cache/src:/go/src \
	    -v ${CWD}/dist:/go/bin \
	    nathanosman/bettergo \
	    go get -pkgdir /go/lib github.com/UnnoTed/fileb0x

cache:
	@mkdir cache

dist:
	@mkdir dist

clean:
	@rm -f server/ab0x.go
	@rm -rf cache dist

.PHONY: clean
