#!/usr/bin/env bash

mkdir -p ./build
go build -o ./build/ ./...

echo "Build folder created"