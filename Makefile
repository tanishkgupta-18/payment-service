.PHONY: run lint test

run:
	go run cmd/app/main.go

lint:
	golangci-lint run

test:
	go test ./...
