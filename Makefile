include local.mk

REQUIRED_BUF_VERSION := latest
REQUIRED_GOLANG_CI_LINT_VERSION := 2.1.6
INSTALLED_GOLANG_CI_LINT_VERSION := $(shell golangci-lint --version 2> /dev/null | grep -Eo '[0-9]+\.[0-9]+\.[0-9]+' | head -1 )
CODE := $(shell find . -type f -name '*.go' -not -name '*mock.go' -not -name '*mock_test.go' -not -name '*_gen.go' -not -name '*.pb.go' -not -path "./vendor/*" -not -path "./output/*")
REQUIRED_JUNITREPORT_VERSION ?= v2.0.0
REQUIRED_GO_IMPORTS_VERSION ?= v0.3.0
GO_IMPORTS := $(shell command -v goimports 2> /dev/null)
REQUIRED_GOFUMPT_VERSION ?= v0.6.0
PKGS := $(shell go list ./...)
CGO_ENABLED ?= 0
APPS = local-chain
PREFIX ?= bin/

$(APPS):
	CGO_ENABLED=$(CGO_ENABLED) go build -mod=vendor -installsuffix cgo -o $(PREFIX)$@ ./cmd/$@

.PHONY: run
run: clear-apps build-apps run-apps

.PHONY: build-apps $(APPS)
build-apps: $(APPS)

.PHONY: clear-apps
clear-apps:
	rm -f bin/$(APPS) && \
	rm -f -r db/*

.PHONY: run-apps
run-apps:
	./bin/$(APPS)

.PHONY: lint
lint:
ifneq "$(INSTALLED_GOLANG_CI_LINT_VERSION)" "$(REQUIRED_GOLANG_CI_LINT_VERSION)"
	CGO_ENABLED=1 go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v$(REQUIRED_GOLANG_CI_LINT_VERSION)
endif
	CGO_ENABLED=1 golangci-lint --timeout 5m run ./...

REQUIRED_MOCKGEN_VERSION ?= v1.6.0
INSTALLED_MOCKGEN_VERSION := $(shell mockgen --version 2> /dev/null)

.PHONY: proto
proto: update lint generate format

.PHONY: update
update: install-buf
	@buf dep update

.PHONY: generate
generate: install-buf clean
ifdef JQ_INSTALLED
		@buf build --exclude-source-info -o -#format=json | jq '.file[] | .name'
else
		@buf build --exclude-source-info
endif
	@buf generate
	@go mod tidy

.PHONY: install-buf
install-buf:
ifndef BUF_INSTALLED
	@go install github.com/bufbuild/buf/cmd/buf@$(REQUIRED_BUF_VERSION)
endif

.PHONY:
clean:
	@rm -fr gen

.PHONY: test
test: unit-test

.PHONY: unit-test
unit-test: install-junit-report install-cobertura
	@go test -short -mod=vendor -cover -count=1 -p=4 -covermode atomic -coverprofile=unit-cover.log $(PWD)/internal \
	-v $(PKGS) 2>&1 | tee unit-test.log
	@go tool cover -func=unit-cover.log | grep -E 'total:\s+\(statements\)\s+'
	@gocover-cobertura < unit-cover.log > unit-coverage.xml
	@cat unit-test.log | go-junit-report -set-exit-code > unit-junit-report.xml

.PHONY: install-junit-report
install-junit-report:
	@go install github.com/jstemmer/go-junit-report/v2@$(REQUIRED_JUNITREPORT_VERSION)

.PHONY: install-cobertura
install-cobertura:
	@go install github.com/boumenot/gocover-cobertura@latest

.PHONY: mock
mock: mockgen format

.PHONY: mockgen
mockgen: install-mockgen
	@go generate -x -run="mockgen" ./...

.PHONY: install-mockgen
install-mockgen:
ifneq "$(INSTALLED_MOCKGEN_VERSION)" "$(REQUIRED_MOCKGEN_VERSION)"
	@echo ======== Installing mockgen ========
	@go install github.com/golang/mock/mockgen@$(REQUIRED_MOCKGEN_VERSION)
endif
	@echo ======== Installed mockgen version ========
	@mockgen --version

.PHONY: format
format: install-gofumpt install-go-imports
	@goimports -local 'local-chain/internal/types' -w $(CODE) && gofumpt -w $(CODE)

.PHONY: install-go-imports
install-go-imports:
ifndef GO_IMPORTS
	@go install golang.org/x/tools/cmd/goimports@$(REQUIRED_GO_IMPORTS_VERSION)
endif

.PHONY: install-gofumpt
install-gofumpt:
	go install mvdan.cc/gofumpt@$(REQUIRED_GOFUMPT_VERSION)

.PHONY: up
up: docker-compose up -d --no-deps --build

.PHONY: down
down: docker-compose down

.PHONY: build
build: go build -v