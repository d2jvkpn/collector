#! /usr/bin/env bash
set -eu -o pipefail
_wd=$(pwd)
_path=$(dirname $0 | xargs -i readlink -f {})


#### install tools
version=23.2

mkdir -p target
wget -P target https://github.com/protocolbuffers/protobuf/releases/download/v${version}/protoc-${version}-linux-x86_64.zip

mkdir -p ~/Apps/protoc-${version}
unzip target/protoc-${version}-linux-x86_64.zip -d ~/Apps/protoc-${version}

####
go get google.golang.org/protobuf@v1.30.0
# go get -u github.com/golang/protobuf/{proto,protoc-gen-go}@v1.30.0
go get -u google.golang.org/grpc

go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

go get google/protobuf/timestamp.proto

#### generate proto
mkdir -p proto

cat > proto/record_data.proto << EOF
syntax = "proto3";
package proto;

option go_package = "github.com/d2jvkpn/collector/proto";

import "google/protobuf/timestamp.proto";

message RecordData {
	string serviceName = 1;
	string serviceVersion = 2;
	string eventId = 3;
	google.protobuf.Timestamp eventAt = 4;
	string bizName = 5;
	string bizVersion = 6;
	map<string, string> bindIds = 7;
	bytes data = 8;
}

message RecordId {
	string id = 1;
}

service DataService {
	rpc Create(Data) returns(RecordId) {};
}
EOF

cat proto/record_data.proto
protoc proto/*.proto --go_out=plugins=grpc:. --go_opt=paths=source_relative --proto_path=.
cat proto/record_data.pb.go

go fmt ./...
go vet ./...

####
# implment...

####
# sed -i '/^\tmustEmbedUnimplemented/s#\t#\t// #' proto/record_data.pb.go
