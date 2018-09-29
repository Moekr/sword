all: build

build:
	./build.sh

fmt:
	find ./ -name "*.go" | grep -v "/vendor/" | xargs gofmt -w

clean:
	rm -rf output