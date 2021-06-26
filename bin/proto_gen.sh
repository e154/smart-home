#!/usr/bin/env bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}")" && cd ../ && pwd)"

cd ${ROOT}/api/protos/

mkdir -p ${ROOT}/api/stub
protoc -I/usr/local/include -I. \
  -I${GOPATH}/src \
  -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
  -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.5.0/protoc-gen-openapiv2 \
  -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0 \
  --grpc-gateway_out=logtostderr=true:${ROOT}/api/stub \
  --openapiv2_out=allow_merge=true,merge_file_name=api,logtostderr=true:${ROOT}/api \
  --go-grpc_out=require_unimplemented_servers=false:${ROOT}/api/stub \
  --go_out=${ROOT}/api/stub \
  *.proto
