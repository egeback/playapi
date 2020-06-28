#!/usr/bin/env bash
#version=0.1.3
version=`cat VERSION`
time=$(date)
swag init --dir ./cmd/playapi/ --parseDependency -output ./internal/docs
go build -o playapi -ldflags="-X 'github.com/egeback/playapi/internal/version.BuildTime=$time' -X 'github.com/egeback/playapi/internal/version.BuildVersion=$version'" ./cmd/playapi/main.go