all:
	go build -o atomicgo .

build:
	go build -o atomicgo .

test:
	go test -v ./...

clean:
	rm atomicgo
