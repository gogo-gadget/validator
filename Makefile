fmt:
	gofmt -w .
make lint:
	golangci-lint run
test:
	go test ./...
