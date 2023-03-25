PROJECT=collective

export CURRENT_DIR=$(shell pwd)

set-env-os-linux:
	export GOOS=linux

set-env-os-win:
	export GOOS=windows

set-env-arch:
	export GOARCH=amd64

clean:
	find ./ -name "*.db" -exec rm -rf {} \;


mod:
	go mod tidy

fmt:
	go fmt ./...

build: fmt set-env-os-linux set-env-arch
	go build

# TODO: Work through race flag exceptions `go test --cover --race ./... `
test: fmt clean
	go test --cover ./... 
	make clean