#!/bin/sh
cd /app
export PATH=$PATH:/usr/local/go/bin
go mod init go-app > /dev/null 2>&1
go get . > /dev/null 2>&1
go run main.go