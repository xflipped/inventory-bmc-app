#!/bin/bash

INCLUDES="-I /usr/local/include -I."

protoFiles="$(find ./ -iname \*.proto)"


rm -rf ./sdk/*

echo "Generating .pb.go files..."

protoc --proto_path=proto $INCLUDES \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	--go_out=./sdk \
	--go-grpc_out=./sdk \
	$protoFiles

echo "Done!"
