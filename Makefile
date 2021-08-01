BIN="./vat-id-validator"
SRC=$(shell find . -name "*.go")

ifeq (, $(shell which golangci-lint))
$(warning "could not find golangci-lint in $(PATH)")
$(info "run: curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh")
endif

.PHONY: fmt lint test clean docker-build

default: all

all: fmt lint test build

build:
	$(info ******************** building recipe stats calculator ********************)
	go build -o vat-id-validator ./cmd

fmt:
	$(info ******************** checking formatting ********************)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

lint:
	$(info ******************** running lint tools ********************)
	golangci-lint run -v

test: lint
	$(info ******************** running tests ********************)
	go test -v -cover ./...

docker-build:
	$(info ******************** building docker image ********************)
	docker build -t vat-id-validator:latest .

clean:
	rm -rf $(BIN)
