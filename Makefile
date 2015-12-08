PACKAGES = $(shell go list ./...) 

all:
	go build -o atomicgo .

build: deps
	go build -o atomicgo .

test:
	go test -v ./...

format:
	@echo "--> Running go fmt" 
	@go fmt $(PACKAGES)

clean:
	rm atomicgo

deps:
	@echo "--> Installing build dependencies"
	@DEP_ARGS="-d -v" sh -c "'$(CURDIR)/scripts/deps.sh'"

updatedeps: deps
	@echo "--> Updating build dependencies"
	@DEP_ARGS="-d -f -u -v" sh -c "'$(CURDIR)/scripts/deps.sh'"
