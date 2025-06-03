GOLANGCI_LINT_VERSION := 2.1.6

all:
	$(MAKE) clean
	$(MAKE) prepare
	$(MAKE) validate

prepare:
	go mod tidy
	go install ./...
	go fmt ./...

validate:
	go vet ./...
	$(MAKE) lint
	$(MAKE) test

lint:
	if [ ! -f ./bin/golangci-lint ]; then \
  		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "./bin" "v${GOLANGCI_LINT_VERSION}"; \
  	fi;
	./bin/golangci-lint run ./...

test:
	CGO_ENABLED=1 go test -race -coverprofile=coverage_out ./...
	go tool cover -func=coverage_out
	go tool cover -html=coverage_out -o coverage.html
	rm -f coverage_out

test_short:
	go test -short ./...

clean:
	go clean