all: build

build:
	./build.sh
	cp -R ./script conf.sample.json sword.*.service ./output

fmt:
	find ./ -name "*.go" | grep -v "/vendor/" | xargs gofmt -w

clean:
	rm -rf output
