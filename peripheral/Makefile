TARGET=itsybitsy-nrf52840

NAME=TBD
VERSION=$(shell cat ../VERSION)-$(shell date +%Y%m%d-%H%M%S)
BRANCH="master"
LDFLAGS="-X main.version=${VERSION} \
		 -X main.name=${NAME}"

.PHONY: build
build:
	tinygo flash \
		-ldflags ${LDFLAGS} \
		-target=$(TARGET) \
		cmd/*.go

.PHONY: tests
tests:
	go test -v ./...
