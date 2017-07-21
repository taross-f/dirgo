#!/bin/bash
go build main.go
GOOS=windows GOARCH=386 go build -o main.exe main.go