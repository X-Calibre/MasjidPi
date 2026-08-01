.PHONY: run fmt tidy test

run:
	cd backend && go run ./cmd/masjidpi

fmt:
	cd backend && gofmt -w .

tidy:
	cd backend && go mod tidy

test:
	cd backend && go test ./...
