VERSION=$(shell cat VERSION)-$(shell date +%Y%m%d-%H%M%S)
BRANCH="master"
LDFLAGS="-X main.version=${VERSION}"

.PHONY: central
central:
	cd central && \
		GOOS=linux GOARCH=arm \
		 go build \
		 -ldflags ${LDFLAGS} \
		 -o gowarm-central \
		 ./cmd/main.go

.PHONY: tests
tests:
	go test -v ./...

