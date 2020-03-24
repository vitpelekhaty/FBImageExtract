GO=go

GOARCH=amd64
GOOS=linux

GOBUILD=${GO} build
GOTEST=${GO} test

TEST_TIMEOUT=-timeout 30s

COMMIT=$(shell git rev-parse --short=8 HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

LDFLAGS=-ldflags "-s -X main.GitBranch=${BRANCH} -X main.GitCommit=${COMMIT}"

BUILD_DIR=./bin

.PHONY: clean test build
all: build

test:
	${GOTEST} ${TEST_TIMEOUT} github.com/vitpelekhaty/fbimgextract/formats/epub
	${GOTEST} ${TEST_TIMEOUT} github.com/vitpelekhaty/fbimgextract/formats/fictionbook/fictionbook2

clean:
	if [ -d "${BUILD_DIR}" ]; then rm -f ${BUILD_DIR}/*; else mkdir ${BUILD_DIR}; fi

build: clean test
	 GOOS=${GOOS} GOARCH=${GOARCH} ${GOBUILD} ${LDFLAGS} -o ${BUILD_DIR}/fbimgextract github.com/vitpelekhaty/fbimgextract
