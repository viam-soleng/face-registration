
selfie-camera: *.go selfiecamera/*.go
	go build -o selfie-camera

test:
	go test

lint:
	gofmt -w -s .

module: selfie-camera
	tar czf module.tar.gz selfie-camera

all: module test