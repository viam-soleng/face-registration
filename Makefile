
selfie-camera:
	go build -o ./bin/selfie-camera

test:
	go test

lint:
	gofmt -w -s .

module: selfie-camera
	tar czf module.tar.gz selfie-camera

all: module test