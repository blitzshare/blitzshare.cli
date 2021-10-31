install:
	#go install golang.org/x/tools/cmd/goimports@latest
	go mod vendor

test:
	ENV=test && go test -v ./... -v -count=1

fix-format:
	gofmt -w -s app/ pkg/ cmd/ mocks/ testhelpers
	goimports -w app/ pkg/ cmd/ mocks/ testhelpers

start:
	go run  app/*.go

build:
	GIN_MODE=release go build -o p2p-client app/main.go

build-docker:
	docker build -t blitzshare.bootstrap.node .

build-docker-run:
	docker build -t blitzshare.bootstrap.node .
	docker run -t blitzshare.bootstrap.node