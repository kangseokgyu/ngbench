#!/bin/bash

echo "=============================="
echo "-- Build for aarch 64"
echo "=============================="

echo "-- Build protobuf"

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/ngbench.proto

mkdir -p build

echo "-- Build reporter_test for aarch 64"

GOOS="linux" \
    GOARCH="arm64" \
    CGO_ENABLED=1 \
    CGO_CFLAGS="-I${PWD}/tmp/libpcap-1.10.5" \
    CGO_LDFLAGS="-L${PWD}/tmp/libpcap-1.10.5 -lpcap" \
    CC="aarch64-unknown-linux-musl-gcc" \
    go test -c github.com/kangseokgyu/ngbench/internal/reporter -o build/reporter.test.arm64

echo "-- Build ngbench-reporter for aarch 64"

GOOS="linux" \
    GOARCH="arm64" \
    CGO_ENABLED=1 \
    CGO_CFLAGS="-I${PWD}/tmp/libpcap-1.10.5" \
    CGO_LDFLAGS="-L${PWD}/tmp/libpcap-1.10.5 -lpcap" \
    CC="aarch64-unknown-linux-musl-gcc" \
    go build -o build/ngbench-reporter.arm64 cmd/ngbench-reporter/ngbench-reporter.go

echo ""
echo "-- Finish to build"
echo "=============================="
