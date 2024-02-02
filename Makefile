BINARY_NAME=gh-devxp

build:
	go build -o ${BINARY_NAME}

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}

install:
	gh extension remove ${BINARY_NAME}
	gh extension install .

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all