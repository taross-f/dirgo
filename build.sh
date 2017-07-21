#!/bin/bash
echo "--build main..."
go build main.go
echo "--build main done!"
echo "--build main.exe..."
GOOS=windows GOARCH=386 go build -o main.exe main.go
echo "--build main.exe done!"
