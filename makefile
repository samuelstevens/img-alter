all: build

build:
	go build main.go

install:
	GOBIN=~/go/bin go install .
