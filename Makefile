LOCAL_BIN:=$(CURDIR)/bin
LOCAL_IN:=$(CURDIR)/internal
GO_KAFKA:=$(LOCAL_IN)/kafka
DB_PORT:=5433

run-prometheus:
	prometheus --config.file config/prometheus.yml

get-mockgen:
	GOBIN=$(LOCAL_BIN) go install github.com/gojuno/minimock/v3/cmd/minimock@latest

sarama-async-prod-mock:
	mkdir -p ${GO_KAFKA}/mocks
	$(LOCAL_BIN)/minimock -i github.com/IBM/sarama.AsyncProducer -o  ${GO_KAFKA}/mocks/async_provider_mock.go

test-kafka-logger:
	go test $(CURDIR)/internal/event_logger/kafka_logger
	go test $(CURDIR)/internal/event_logger/kafka_logger -coverprofile=$(LOCAL_BIN)/coverage1
	go tool cover -html $(LOCAL_BIN)/coverage1 -o $(LOCAL_BIN)/index.html


get-protoc:
	apt install -y protobuf-compiler
	protoc --version

all: deps generate build run

deps: .vendor-proto
	docker-compose -f docker/docker-compose.yml up -d
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@latest

generate:
	protoc --proto_path api --proto_path vendor.protogen \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go --go_out=${LOCAL_IN} --go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc --go-grpc_out=${LOCAL_IN} --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway --grpc-gateway_out ${LOCAL_IN} --grpc-gateway_opt paths=source_relative \
		--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 --openapiv2_out=${LOCAL_IN} \
		--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate --validate_out="lang=go,paths=source_relative:${LOCAL_IN}" \
		./api/pvz-service/v1/pvz_service.proto

build:
	$(LOCAL_BIN)/goose -dir $(CURDIR)/migrations postgres "postgres://postgres:postgres@localhost:$(DB_PORT)/postgres?sslmode=disable" up
	touch  $(LOCAL_BIN)
	go build -o $(LOCAL_BIN)/pvz-service cmd/pvz_service/main.go
	go build -o $(LOCAL_BIN)/pvz cmd/pvz_cli/main.go

run:
	docker-compose -f docker/docker-compose.yml start
	touch logs.txt
	$(LOCAL_BIN)/pvz-service > logs.txt 2>&1 &
	sleep 1
	$(LOCAL_BIN)/pvz inter
	pkill -SIGINT pvz-service
	docker-compose -f docker/docker-compose.yml stop


.vendor-proto: .vendor-proto/google/protobuf .vendor-proto/google/api .vendor-proto/protoc-gen-openapiv2/options .vendor-proto/validate

.vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/grpc-ecosystem && \
 		cd vendor.protogen/grpc-ecosystem && \
		git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
		git checkout
		mkdir -p vendor.protogen/protoc-gen-openapiv2
		mv vendor.protogen/grpc-ecosystem/protoc-gen-openapiv2/options vendor.protogen/protoc-gen-openapiv2
		rm -rf vendor.protogen/grpc-ecosystem

.vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor.protogen/protobuf &&\
		cd vendor.protogen/protobuf &&\
		git sparse-checkout set --no-cone src/google/protobuf &&\
		git checkout
		mkdir -p vendor.protogen/google
		mv vendor.protogen/protobuf/src/google/protobuf vendor.protogen/google
		rm -rf vendor.protogen/protobuf

.vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor.protogen/googleapis && \
 		cd vendor.protogen/googleapis && \
		git sparse-checkout set --no-cone google/api && \
		git checkout
		mkdir -p  vendor.protogen/google
		mv vendor.protogen/googleapis/google/api vendor.protogen/google
		rm -rf vendor.protogen/googleapis

.vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor.protogen/tmp && \
		cd vendor.protogen/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor.protogen/validate
		mv vendor.protogen/tmp/validate vendor.protogen/
		rm -rf vendor.protogen/tmp

