#!/bin/bash

cd `dirname $0`

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/win/AndreFileUploader.exe main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/linux/AndreFileUploader main.go


