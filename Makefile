
all: test

deps:
	go get -d -v ./...

updatedeps:
	go get -d -v -u -f ./...

testdeps: deps
	go get -d -v -t ./...

updatetestdeps: deps
	go get -d -v -t -u -f ./...

build: deps
	go build ./...

test: testdeps
	go test -test.v -race ./...

cov: testdeps
	go get -v github.com/axw/gocov/gocov
	gocov test > tmp/coverage.json
	gocov annotate -color -ceiling 100 tmp/coverage.json

doc:
	go get -v github.com/robertkrimen/godocdown/godocdown
	cp .readme.header README.md
	godocdown | tail -n +7 >> README.md
