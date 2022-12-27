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

test: fmt clean
	go test --cover ./... 