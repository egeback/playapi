#!/usr/bin/env bash
version=2
time=$(date)
go build -o main -ldflags="-X 'main.BuildTime=$time' -X 'main.BuildVersion=$version'" ./internal/main.go