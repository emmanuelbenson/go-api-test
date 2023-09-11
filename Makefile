build:
	@go build -o bin/api-test

run: build
	@./bin/api-test

test:
	@go test -v ./...