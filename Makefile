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
	go run cmd/main.go --init

start-connect:
	go run cmd/main.go --connect

build:
	go build -o p2p-client cmd/main.go

build-docker:
	docker build -t blitzshare.bootstrap.node .

build-docker-run:
	docker build -t blitzshare.bootstrap.node .
	docker run -t blitzshare.bootstrap.node

build-mocks:
	mockery --all --dir "./app/"