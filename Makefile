VERSION=1.0.0
GITHASH ?= $(shell git describe --match 'v[0-9]*' --tags --always)

build_package:
	@echo "Building Genc binary..."
	env CGO_ENABLED=0 go build -ldflags "-X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash=$(GITHASH) -X main.version=${VERSION}" -o gop3p cmd/main.go
	@echo "Done build!"

install:
	cp gop3p /usr/local/bin/
	@echo "Done install!"

all: build_package install
