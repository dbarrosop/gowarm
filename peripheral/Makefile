include ../build/makefiles/common.mk
include ../build/makefiles/go.mk

TARGET=itsybitsy-nrf52840

.PHONY: build
build:  ## Build project and flash it
	tinygo flash \
		-ldflags ${LDFLAGS} \
		-target=$(TARGET) \
		cmd/*.go

.PHONY: tests
tests: go-tests go-lint ## Run tests
