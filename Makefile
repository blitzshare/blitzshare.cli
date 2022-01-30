SHELL := /bin/bash
CWD := $(shell cd -P -- '$(shell dirname -- "$0")' && pwd -P)

export GO111MODULE := on
export GOBIN := $(CWD)/.bin

install:
	go install $(shell go list -f '{{join .Imports " "}}' tools.go)
	go mod vendor
test:
	go test  --tags='test' -v ./app/... -v -count=1 -cover -coverprofile=coverage.out
coverage-report-html:
	go tool cover -html=coverage.out	
fix-format:
	gofmt -w -s app/ pkg/ cmd/ mocks/ testhelpers
	goimports -w app/ pkg/ cmd/ mocks/ testhelpers
start:
	go run cmd/main.go --start
start-fs:
	go run cmd/main.go --start --file ./README.md
start-connect:
	go run cmd/main.go --connect
build-mocks:
	mockery --all --dir "./app/"
build:
	go build ./cmd/main.go