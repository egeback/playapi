#!/usr/bin/env bash
#version=0.1.3
version=`cat VERSION`
time=$(date)
go build -o main -ldflags="-X 'github.com/egeback/play_media_api/internal/version.BuildTime=$time' -X 'github.com/egeback/play_media_api/internal/version.BuildVersion=$version'" ./internal/main.go