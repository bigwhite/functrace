.PHONY: all build install clean

all: build

build:
	go build github.com/bigwhite/functrace/cmd/gen

install: 
	go install github.com/bigwhite/functrace/cmd/gen

clean:
	go clean
	rm -fr gen
