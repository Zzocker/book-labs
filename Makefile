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

proto:
	protoc -I protos/ --go_out ${GOPATH}/src/ protos/common/common.proto
	protoc -I protos/ --go_out ${GOPATH}/src/ protos/mediafile/mediafile.proto
	protoc -I protos/ --go-grpc_out ${GOPATH}/src/ protos/mediafile/mediafile.proto

demo:
	docker-compose -f docker-compose.yaml up -d

migrate:
	./migration/s3.sh