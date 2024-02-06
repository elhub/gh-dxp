BINARY_NAME=gh-devxp
BUILD_DIR=build

build:
	go build -o ${BUILD_DIR}/${BINARY_NAME}

run: build
	${BUILD_DIR}/${BINARY_NAME}

clean:
	go clean
	rm -rf ${BUILD_DIR}

install:
	gh extension remove ${BINARY_NAME}
	gh extension install .

test:
	mkdir -p ${BUILD_DIR}
	go test ./... -coverprofile=${BUILD_DIR}/coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all

ci_test:
	cd .teamcity && mvn compile && cd ..