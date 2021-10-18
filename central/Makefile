include ../build/makefiles/common.mk
include ../build/makefiles/go.mk

NAME=central

.PHONY: build
build:
	GOOS=linux GOARCH=arm \
	 go build \
	 -ldflags ${LDFLAGS} \
	 -o gowarm-central \
	 ./cmd/main.go

.PHONY: tests
tests:
	go test -v ./...

