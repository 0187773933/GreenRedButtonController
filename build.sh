#!/bin/bash
rm -rf ./bin
GOOS=linux GOARCH=arm64 go build -o bin/linux/arm64/gr
