BINARY_NAME=go-fast-cdn
OS_NAME := $(shell uname -s | tr A-Z a-z)
# Detect if we're running in MINGW64 and treat it as Windows
ifeq ($(findstring mingw,$(OS_NAME)),mingw)
    OS_NAME := windows
endif
ARCH := $(shell uname -m | tr A-Z a-z | sed 's/^aarch/arm/' | sed 's/^x86_64/amd64/')

prep:
	go mod tidy
	go mod download
	cd ui && pnpm i
	go install github.com/air-verse/air@latest

build: build_ui build_bin

build_ui:
	pnpm --dir ./ui build

build_bin:
ifeq ($(OS_NAME),darwin)
	GOARCH=${ARCH} GOOS=darwin CGO_ENABLED=0 go build -o bin/${BINARY_NAME}-darwin
else ifeq ($(OS_NAME),linux)
	CC="x86_64-linux-musl-gcc" GOARCH=${ARCH} GOOS=${OS_NAME} CGO_ENABLED=0 go build -o bin/${BINARY_NAME}-${OS_NAME}
else ifeq ($(OS_NAME),windows)
	CC="x86_64-w64-mingw32-gcc" GOARCH=${ARCH} GOOS=windows CGO_ENABLED=0 go build -o bin/${BINARY_NAME}-windows.exe
endif

run: build
ifeq ($(OS_NAME),windows)
	bin/${BINARY_NAME}-windows.exe
else
	bin/${BINARY_NAME}-${OS_NAME}
endif

dev:
	air

clean: 
	go clean
	rm -rf bin/*
	rm -rf ui/build/*

vet:
	go vet

test:
	go test ./...
