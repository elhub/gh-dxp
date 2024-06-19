BINARY_NAME=gh-dxp
BUILD_DIR=build

.PHONY: all clean check

all: clean dep check vet lint build

build:
	go build -o ${BUILD_DIR}/${BINARY_NAME}

check:
	mkdir -p ${BUILD_DIR}
	go test ./... -coverprofile=${BUILD_DIR}/coverage.out

clean:
	go clean
	rm -rf ${BUILD_DIR}
	rm -rf ${BINARY_NAME}

dep:
	go mod download

install: clean build
	-gh extension remove ${BINARY_NAME}
	cp ${BUILD_DIR}/${BINARY_NAME} .; gh extension install .

run: build
	${BUILD_DIR}/${BINARY_NAME}

vet:
	go vet

teamcityCheck:
	cd .teamcity && mvn teamcity-configs:generate
