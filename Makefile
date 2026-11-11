.PHONY: build test test-integration lint fmt mocks docker help

build:
	go build ./...

test:
	go test -race $$(go list ./tests/... | grep '_unit_tests')

# Needs Docker: starts a real Kubo through testcontainers.
test-integration:
	go test -race $$(go list ./tests/... | grep '_integration_tests')

lint:
	golangci-lint run ./...

fmt:
	golangci-lint fmt ./...

# Needs deepsearch-mockgen-cli and a TTY. mockgen duplicates the generated client's interfaces; keep one file.
mocks:
	deepsearch-mockgen-cli -S ./src -O ./tests -A -P
	rm -f tests/kubo_mcp_http_kubo_mocks/clients_mocks/client_mock.go \
		tests/kubo_mcp_http_kubo_mocks/clients_mocks/http_request_doer_mock.go

docker:
	docker build -f src/kubo_mcp_host/Dockerfile -t kubo-mcp:dev .

help:
	@echo "build  - go build ./..."
	@echo "test   - go test -race on unit-test packages"
	@echo "test-integration - tests against a real Kubo (needs Docker)"
	@echo "lint   - golangci-lint run (needs golangci-lint v2)"
	@echo "fmt    - golangci-lint fmt, 120-char cap"
	@echo "mocks  - regenerate mocks (needs deepsearch-mockgen-cli)"
	@echo "docker - build the image as kubo-mcp:dev"
