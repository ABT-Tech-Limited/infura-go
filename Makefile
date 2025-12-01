.PHONY: test

tidy:
	@echo "Running tidy..."
	@go mod tidy
	@goimports -w .
	@go vet ./...

test:
	@echo "Running tests..."
	@go mod tidy
	@goimports -w .
	@go vet ./...
	@GOMAXPROCS=1 go test -p=1 ./... -v