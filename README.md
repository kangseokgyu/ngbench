# ngbench

## Update proto

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/ngbench.proto
```

## Test

```bash
go test ./...
```

## libpcap

pcap을 사용하려면 libpcap이 설치되어 있어야 한다.

```bash
brew install libpcap
```

pcap을 사용하면서 크로스 컴파일을 하기 위해서는 타겟에 맞는 pcap 라이브러리가 빌드되어 있어야 한다.

### Install Toolchain

```bash
brew tap messense/macos-cross-toolchains
brew install aarch64-unknown-linux-musl
```

### Cross compile

```bash
wget https://www.tcpdump.org/release/libpcap-1.10.5.tar.xz
tar -zxvf libpcap-1.10.5.tar.xz
cd libpcap-1.10.5
export CC=aarch64-linux-musl-gcc
export CCFLAGS=-I/opt/homebrew/Cellar/aarch64-unknown-linux-musl/13.3.0/toolchain/aarch64-unknown-linux-musl/sysroot/usr/include
./configure --host=aarch64-unknown-linux-musl --with-pcap=linux
make
ln -s libpcap.so.1.10.5 libpcap.so
```

### Link libpcap

```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CGO_CFLAGS="-I${PWD}/tmp/libpcap-1.10.5" CGO_LDFLAGS="-L${PWD}/tmp/libpcap-1.10.5 -lpcap" go test -c github.com/kangseokgyu/ngbench/internal/reporter -o reporter.test.arm64
```
