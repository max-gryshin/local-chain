REQUIRED_GOLANG_CI_LINT_VERSION := 2.1.6
INSTALLED_GOLANG_CI_LINT_VERSION := $(shell golangci-lint --version 2> /dev/null | grep -Eo '[0-9]+\.[0-9]+\.[0-9]+' | head -1 )

REQUIRED_JUNITREPORT_VERSION ?= v2.0.0

.PHONY: lint
lint:
ifneq "$(INSTALLED_GOLANG_CI_LINT_VERSION)" "$(REQUIRED_GOLANG_CI_LINT_VERSION)"
	CGO_ENABLED=1 go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v$(REQUIRED_GOLANG_CI_LINT_VERSION)
endif
	CGO_ENABLED=1 golangci-lint --timeout 5m run ./...

.PHONY: test
test: unit-test

.PHONY: unit-test
unit-test: install-junit-report install-cobertura
	@go test -short -mod=vendor -cover -count=1 -p=4 -covermode atomic -coverprofile=unit-cover.log \
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

.PHONY: up
up: docker-compose up -d --no-deps --build

.PHONY: down
down: docker-compose down

.PHONY: build
build: go build -v