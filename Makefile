# ==============================================================================
# Project

test:
	go test ./... -count=1

coverage:
		go test ./... -count=1 -coverprofile /tmp/cover.out
		go tool cover -html=/tmp/cover.out

lint:
	golangci-lint run

fix-lint:
	golangci-lint run --fix

tidy:
	go mod tidy

proto:
	protoc -I ./protos --go_out ${GOPATH}/src/ protos/common/common.proto
	protoc -I ./protos --go_out ${GOPATH}/src/ protos/profile/profile.proto
	protoc -I ./protos --go-grpc_out ${GOPATH}/src/ protos/profile/profile.proto
	protoc -I ./protos --go_out ${GOPATH}/src/ protos/book/book.proto
	protoc -I ./protos --go-grpc_out ${GOPATH}/src/ protos/book/book.proto