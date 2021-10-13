GOLANGCI_LINT_VER="v1.41.1"
GOTEST_OPTIONS?=-v ./... -coverprofile=coverage.txt -covermode=atomic

.PHONY: go-tests
go-tests: ## Run go test
	go test $(GOTEST_OPTIONS)


.PHONY: go-integration-tests
go-integration-tests: ## Run go test with integration flags
	go test -tags=integration $(GOTEST_OPTIONS) -integration

.PHONY: go-install-linter
go-install-linter: ## Install golangci-linet
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(GOLANGCI_LINT_VER)
	mkdir -p $(BASE_PATH)/bin
	mv bin/golangci-lint $(BASE_PATH)/bin/
	rm -rf ./bin

.PHONY: go-lint
go-lint: ## Run Go linters in a Docker container
	../bin/golangci-lint \
		run \
		--timeout 300s
