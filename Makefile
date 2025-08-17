build:
	@go build -o bin/ypoker

run: build
	@./bin/ypoker

dev:
	@air

test:
	go test ./...