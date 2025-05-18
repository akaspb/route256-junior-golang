wsl --set-default Ubuntu
tar -xvf go1.23.2.linux-amd64.tar.gz
export PATH=$PATH:/mnt/d/Go/go/bin
kill -9 1110

..\prometheus\prometheus.exe --config.file "configs/prometheus.yml"

# echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# go build .\cmd\main.go
# .\main.exe inter

go run .\cmd\main.go inter

go build -o pvz cmd/main.go
./pvz inter
.\pvz.exe inter -s="08.09.2024"

..\goose.exe -dir ./migrations postgres "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable" up

receive -o=2 -m=200 -w=5 -c=1 -e=12.10.2024
receive -o=2 -m=200 -w=5 -c=1 -e=27.10.2024
receive -o=3 -m=500 -w=5 -c=1 -e=27.10.2024 -p=wrap
receive -o=4 -m=500 -w=9 -c=1 -e=27.10.2024 -p=packet
receive -o=5 -m=700 -w=29 -c=1 -e=27.10.2024 -p=box
receive -o=5 -m=700 -w=29 -c=1 -e=27.10.2024 -p=box
receive -o=9 -m=250 -w=77.7 -c=2 -e=27.10.2024 -p=wrap

receive -o=8 -m=700 -w=29 -c=2 -e=27.10.2024 -p=box
give 1 8
give 2 9

give 1 2 3 4 5 7
give 1 2 3 4 5
list 1
give 2 8
return -o=2 -c=2
return -o=2 -c=1
return -o=3 -c=1
return -o=4 -c=1
return -o=5 -c=1
return -o=7 -c=2
return -o=8 -c=2
return -o=9 -c=2
give 1 2 3 4 5
returns -o=0 -l=15
remove 2
remove 3
remove 4
remove 5
remove 7
remove 8
remove 9


..\minimock_3.4.0_windows_amd64\minimock.exe ./internal/storage/storage.go

# # https://www.technewstoday.com/install-and-use-make-in-windows/
# # Winget install GnuWin32.make

# # https://gnuwin32.sourceforge.net/install.html
# ..\make-3.81-bin\bin\make.exe

winget install --id=GnuWin32.Make  -e

go test ./test/cli -coverprofile=coverage
go tool cover -html coverage -o index.html


docker-compose -f docker/docker-compose.yml up -d
docker-compose -f docker/docker-compose.yml start
docker-compose -f docker/docker-compose.yml stop
docker-compose -f docker/docker-compose.yml restart
docker-compose -f docker/docker-compose.yml ps

docker exec -it route256_db psql -U postgres -d postgres
=# \?
=# \d

docker exec -it route256-init-topics /bin/bash
cd /opt/kafka/kafka/
./bin/kafka-storage.sh random-uuid

# goose DRIVER DBSTRING create NAME TYPE
..\goose.exe -dir ./migrations postgres "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable" status

..\goose.exe -dir ./migrations postgres "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable" create add_packaging sql
..\goose.exe -dir ./migrations postgres "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable" create add_statuses sql
..\goose.exe -dir ./migrations postgres "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable" create add_orders sql
..\goose.exe -dir ./migrations postgres "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable" create add_packs sql
..\goose.exe -dir ./migrations postgres "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable" create create_index_packaging sql
..\goose.exe -dir ./migrations postgres "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable" create create_index_orders sql

..\goose.exe -dir ./migrations postgres "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable" up


../protoc.exe --version

go env
go env -w GOBIN="D:\Go\Ozon\Homework\bin\"
go env -w GOBIN=/mnt/d/Go/Ozon/Homework
go env -u GOBIN

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go install github.com/envoyproxy/protoc-gen-validate@latest
go install github.com/rakyll/statik@latest

git clone -b main --single-branch -n --depth=1 --filter=tree:0 `
    https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/grpc-ecosystem
cd vendor.protogen/grpc-ecosystem
git sparse-checkout set --no-cone protoc-gen-openapiv2/options
git checkout
mkdir -p vendor.protogen/protoc-gen-openapiv2
mv vendor.protogen/grpc-ecosystem/protoc-gen-openapiv2/options vendor.protogen/protoc-gen-openapiv2
rm -rf vendor.protogen/grpc-ecosystem

git clone -b main --single-branch -n --depth=1 --filter=tree:0 `
    https://github.com/protocolbuffers/protobuf vendor.protogen/protobuf
cd vendor.protogen/protobuf
git sparse-checkout set --no-cone src/google/protobuf
git checkout
mv vendor.protogen/protobuf/src/google/protobuf vendor.protogen/google
rm -rf vendor.protogen/protobuf

git clone -b master --single-branch -n --depth=1 --filter=tree:0 `
    https://github.com/googleapis/googleapis vendor.protogen/googleapis
cd vendor.protogen/googleapis
git sparse-checkout set --no-cone google/api
git checkout
mv vendor.protogen/googleapis/google/api vendor.protogen/google
rm -rf vendor.protogen/googleapis

git clone -b main --single-branch --depth=2 --filter=tree:0 `
    https://github.com/bufbuild/protoc-gen-validate vendor.protogen/tmp
cd vendor.protogen/tmp
git sparse-checkout set --no-cone validate
git checkout
mv vendor.protogen/tmp/validate vendor.protogen/
rm -rf vendor.protogen/tmp

../protoc --proto_path api --proto_path vendor.protogen `
    --plugin=protoc-gen-go="./bin/protoc-gen-go.exe" `
    --go_out="./pkg/" `
    --go_opt=paths=source_relative `
    --plugin=protoc-gen-go-grpc="./bin/protoc-gen-go-grpc.exe" --go-grpc_out="./pkg/" --go-grpc_opt=paths=source_relative `
    --plugin=protoc-gen-grpc-gateway="./bin/protoc-gen-grpc-gateway.exe" --grpc-gateway_out ./pkg/ --grpc-gateway_opt paths=source_relative `
    --plugin=protoc-gen-openapiv2="./bin/protoc-gen-openapiv2.exe" --openapiv2_out=./pkg/ `
    --plugin=protoc-gen-validate="./bin/protoc-gen-validate.exe" --validate_out="lang=go,paths=source_relative:./pkg/" `
    ./api/pvz-service/v1/pvz_service.proto

# ../protoc --go_out=plugins=grpc:. ./api/pvz-service/v1/pvz_service.proto


http://localhost:8080

if err != nil {
    log.Fatal(err)
}


https://prometheus.io/docs/visualization/grafana/#installing
