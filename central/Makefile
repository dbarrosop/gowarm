include ../build/makefiles/common.mk
include ../build/makefiles/go.mk


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

