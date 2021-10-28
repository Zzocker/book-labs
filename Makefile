# ====================================================
# Project

tidy:
	go mod tidy

test:
	go test ./... -count=1

coverage:
	go test ./... -count=1 -coverprofile /tmp/cover.out
	go tool cover -html=/tmp/cover.out
lint:
	golangci-lint run

fix-lint:
	golangci-lint run --fix