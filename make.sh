#!/usr/bin/env sh

VERSION=$(cat VERSION)
NAME=mqtt-tempo

rm -rf release/$VERSION

echo build release/$VERSION/linux/arm
mkdir -p release/$VERSION/linux/arm
GOOS=linux GOARCH=arm GOARM=7 go build $NAME.go
mv $NAME release/$VERSION/linux/arm

echo build release/$VERSION/linux/amd64
mkdir -p release/$VERSION/linux/amd64
GOOS=linux GOARCH=amd64 go build $NAME.go
mv $NAME release/$VERSION/linux/amd64

echo build release/$VERSION/windows/amd64
mkdir -p release/$VERSION/windows/amd64
GOOS=windows GOARCH=amd64 go build $NAME.go
mv $NAME.exe release/$VERSION/windows/amd64

echo build release/$VERSION/linux/aarch64
mkdir -p release/$VERSION/linux/aarch64
GOOS=linux GOARCH=arm64 go build $NAME.go
mv $NAME release/$VERSION/linux/aarch64
