#!/bin/bash
go mod tidy
go build -o main .

zip archive.zip main
