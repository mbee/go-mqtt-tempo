#!/usr/bin/env sh

VERSION=1.1.1
PRODUCT=mqtt-tempo

echo fetch dependencies
go get github.com/Sirupsen/logrus
go get github.com/eclipse/paho.mqtt.golang
go get golang.org/x/sys/unix
go get github.com/konsorten/go-windows-terminal-sequences

echo build linux/arm/5
mkdir -p release/$VERSION/linux/arm
GOOS=linux GOARCH=arm GOARM=5 go build $PRODUCT.go
mv $PRODUCT release/$VERSION/linux/arm/$PRODUCT

echo build linux/amd64
mkdir -p release/$VERSION/linux/amd64
GOOS=linux GOARCH=amd64 go build $PRODUCT.go
mv $PRODUCT release/$VERSION/linux/amd64/$PRODUCT

echo build windows/amd64
mkdir -p release/$VERSION/windows/amd64
GOOS=windows GOARCH=amd64 go build $PRODUCT.go
mv $PRODUCT.exe release/$VERSION/windows/amd64/$PRODUCT
