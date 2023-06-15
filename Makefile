
all: build

build:
	go build github.com/pengxuan37/functrace/cmd/gen

install: 
	go install github.com/pengxuan37/functrace/cmd/gen
clean:
	go clean
	rm -fr gen
