build:
	@go build -o bin/ypoker

run: build
	@./bin/ypoker

test:
	go test ./...