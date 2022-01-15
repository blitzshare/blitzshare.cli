install:
	go install golang.org/x/tools/cmd/goimports@latest
	go get -d github.com/vektra/mockery/v2/.../
	go mod vendor

test:
	ENV=test && go test -v ./... -v -count=1 && echo $?

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