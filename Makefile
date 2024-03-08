GOBASE = $(shell pwd)

LINT_PATH = $(GOBASE)/build/lint

deps: ## Fetch required dependencies
	go mod tidy -compat=1.20
	go mod download

lint: install-golangci ## Linter for developers
	$(LINT_PATH)/golangci-lint run --timeout=5m --build-tags=$(goKafkaTag) -c .golangci.yml

install-golangci: ## Install the correct version of lint
	@if [ ! -f $(LINT_PATH)/golangci-lint ] ; then GOBIN=$(LINT_PATH) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.2 ; fi

test:
	go test ./...
