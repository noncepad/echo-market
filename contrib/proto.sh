#!/bin/bash

set -e

decho() {
  1>&2 echo $@
}

PROTO_BIN="protoc"

PROTO_DIR="./contrib/market/proto"
PROTO_GO_DIR="./proto"

test_all() {
  decho $@
}

build_go() {
  rm -r $PROTO_GO_DIR 2>>/dev/null || true
  $PROTO_BIN --experimental_allow_proto3_optional --proto_path=${PROTO_DIR} --go-grpc_out=. --go_out=. $@
  mv "github.com/noncepad/echo-market/proto" ./
  rm -r github.com
}

build_go $(echo $(find ${PROTO_DIR} -type f))

