GO=go

.PHONY: all build test clean

all: build 

build:
	cd ./cmd/gopaths && ${GO} build 

test:
	${GO} test ./...

clean:
	rm ./cmd/gopaths/gopaths
