BASE_PATH=$(shell dirname $(abspath $(lastword $(MAKEFILE_LIST)/../..)))
FULL_PATH=$(shell dirname $(abspath $(firstword $(MAKEFILE_LIST))))
REL_PATH=$(subst $(BASE_PATH)/,,$(FULL_PATH))

BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
VERSION=$(shell cat $(BASE_PATH)/VERSION)+$(BRANCH).$(shell date +%Y%m%d-%H%M%S)

DEVICE_NAME=TBD
LDFLAGS="-X main.version=${VERSION} \
		 -X main.name=${DEVICE_NAME}"

PROJECT=github.com/dbarrosop/gowarm
NAME?=$(shell basename $(FULL_PATH))

.PHONY: info
info: ## Echo common env vars
	@echo BASE_PATH: $(BASE_PATH)
	@echo FULL_PATH: $(FULL_PATH)
	@echo REL_PATH: $(REL_PATH)
	@echo PROJECT: $(PROJECT)
	@echo NAME: $(NAME)
	@echo VERSION: $(VERSION)
	@echo DEVICE_NAME: $(DEVICE_NAME)

help: ## Show this help.
	@IFS=$$'\n' ; \
	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//'`); \
	for help_line in $${help_lines[@]}; do \
			IFS=$$'#' ; \
			help_split=($$help_line) ; \
			help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
			help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
			printf "%-30s %s\n" $$help_command $$help_info ; \
	done
