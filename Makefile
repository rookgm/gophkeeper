
APPNAME := gophkeeper
VERSION := $(shell echo "v1.0.0")
BUILD_DATE := $(shell date +'%Y/%m/%d %H:%M:%S')
COMMIT := $(shell echo "#1")

MAIN_PACKAGE_PATH := main

LDFLAGS := -ldflags "-w -s \
	-X $(MAIN_PACKAGE_PATH).BuildVersion=$(VERSION) \
	-X '$(MAIN_PACKAGE_PATH).BuildDate=$(BUILD_DATE)' \
	-X $(MAIN_PACKAGE_PATH).BuildCommit=$(COMMIT)"

.PHONY: build
build: server client

.PHONY: client
client:
	go build -buildvcs=false $(LDFLAGS) -o bin/$(APPNAME)-client ./cmd/client

.PHONY: server
server:
	go build -buildvcs=false $(LDFLAGS) -o bin/$(APPNAME)-server ./cmd/server

.PHONY: all
all:
	#client amd64
	# linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false $(LDFLAGS) -o bin/$(APPNAME)-client-linux-amd64 ./cmd/client
	# windows
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false $(LDFLAGS) -o bin/$(APPNAME)-client-windows-amd64.exe ./cmd/client
	# macos
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false $(LDFLAGS) -o bin/$(APPNAME)-client-darwin-amd64 ./cmd/client

	# server amd64
	# linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false $(LDFLAGS) -o bin/$(APPNAME)-server-linux-amd64 ./cmd/server
	# windows
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false $(LDFLAGS) -o bin/$(APPNAME)-server-windows-amd64.exe ./cmd/server
	# macos
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -buildvcs=false $(LDFLAGS) -o bin/$(APPNAME)-server-darwin-amd64 ./cmd/server

PHONY: clean
clean:
	rm -rf bin/*