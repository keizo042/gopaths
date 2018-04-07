GO=go

.PHONY: all build test clean

all: build 

build:
	${GO} build ./cmd/gopaths 

test:
	${GO} test ./...

clean:
	rm ./cmd/gopaths/gopaths
