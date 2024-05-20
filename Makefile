lint:
	golangci-lint run ./...

test:
	@go test -race -covermode=atomic -coverprofile=coverage.out ./...

run:
	@go run cmd/gophermart/main.go
