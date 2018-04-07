GO=go

.PHONY: all build test clean

all: build 

build:
	${GO} build ./cmd/gopaths/ -o /cmd/gopaths/gopaths

test:
	${GO} test ./...

clean:
	rm ./cmd/gopaths/gopaths
