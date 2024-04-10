
face-camera: *.go facecamera/*.go
	go build -o face-camera

test:
	go test

lint:
	gofmt -w -s .

module: face-camera
	tar czf module.tar.gz face-camera

all: module test