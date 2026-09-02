# Makefile
.PHONY: build run clean dev

build:
	@echo "Building binary..."
	@go build -o bin/drive-svc cmd/server/main.go

run: build
	@echo "Running..."
	@./bin/drive-svc

clean:
	@rm -rf bin/

dev:
	@echo "Starting dev server..."
	@/home/sirkartik/go/bin/air
