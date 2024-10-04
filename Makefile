BINARY_NAME=gh-dxp
BUILD_DIR=build
BIN_DIR=${BUILD_DIR}/gh-dxp

.PHONY: all clean check

all: clean dep check vet build

build:
	go build -o ${BIN_DIR}/${BINARY_NAME}

check:
	mkdir -p ${BUILD_DIR}
	go test ./... -coverprofile=${BUILD_DIR}/coverage.out
	go tool cover -func build/coverage.out | grep "total"

clean:
	go clean
	rm -rf ${BUILD_DIR}
	rm -rf ${BINARY_NAME}

dep:
	go mod download

install: clean build
	-gh extension remove ${BINARY_NAME}
	cd ${BIN_DIR}; gh extension install .

run: build
	${BIN_DIR}/${BINARY_NAME}

vet:
	go vet

teamcityCheck:
	cd .teamcity && mvn teamcity-configs:generate
