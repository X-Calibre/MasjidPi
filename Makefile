.PHONY: build run fmt tidy test

build:
	cd backend && mkdir -p build && go build -o build/masjidpi ./cmd/masjidpi

run:
	cd backend && go run ./cmd/masjidpi

fmt:
	cd backend && gofmt -w .

tidy:
	cd backend && go mod tidy

test:
	cd backend && go test ./...
