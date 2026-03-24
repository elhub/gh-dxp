BINARY_NAME=gh-dxp
BUILD_DIR=build
BIN_DIR=${BUILD_DIR}/gh-dxp
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null)
REPO ?= $(shell git config --get remote.origin.url | sed -E 's#(git@github.com:|https://github.com/)##; s#\.git$$##')

PLATFORMS=linux/amd64 darwin/arm64

# Phony targets
.PHONY: all build check clean dep install owasp run vet release teamcity-check

# Default target
all: clean dep check vet build

# Target: Build the binary
build:
	go build -o ${BIN_DIR}/${BINARY_NAME}

# Target: Run tests and generate coverage report
check:
	mkdir -p ${BUILD_DIR}
	go test ./... -coverprofile=${BUILD_DIR}/coverage.out
	go tool cover -html=${BUILD_DIR}/coverage.out -o ${BUILD_DIR}/coverage.html
	go tool cover -func ${BUILD_DIR}/coverage.out | grep "total"

# Target: Clean build artifacts
clean:
	go clean
	rm -rf ${BUILD_DIR}
	rm -rf ${BINARY_NAME}

# Target: Download dependencies
dep:
	go mod download

# Target: Install the binary as a GitHub CLI extension
install: clean build
	-gh extension remove ${BINARY_NAME}
	cd ${BIN_DIR}; gh extension install .

# Target: Run OWASP security checks
owasp:
	go run github.com/securego/gosec/v2/cmd/gosec ./...

# Target: Run the binary
run: build
	${BIN_DIR}/${BINARY_NAME}

# Target: Run Go vet for static analysis
vet:
	go vet

# Target: Build release binaries for all platforms
release:
	@test -n "${VERSION}" || (echo "No git tags found. Create a tag before running make release."; exit 1)
	@echo "Preparing release for ${VERSION}"
	@if ! gh release view "${VERSION}" --repo "${REPO}" >/dev/null 2>&1; then \
		mkdir -p dist; \
		for platform in $(PLATFORMS); do \
			GOOS=$$(echo $$platform | cut -d/ -f1); \
			GOARCH=$$(echo $$platform | cut -d/ -f2); \
			echo "Building for $$GOOS/$$GOARCH..."; \
			GOOS=$$GOOS GOARCH=$$GOARCH go build -ldflags="-X 'main.version=$(VERSION)'" -o "dist/$(BINARY_NAME)-$$GOOS-$$GOARCH" . || exit 1; \
		done; \
		gh release create "${VERSION}" dist/$(BINARY_NAME)-* --repo "${REPO}" --verify-tag --generate-notes; \
		echo "All builds completed successfully."; \
	else \
		echo "Release ${VERSION} already exists in GitHub. Refusing to publish."; \
		exit 1; \
	fi

# Target: Generate TeamCity configuration (for testing purposes)
teamcity-check:
	cd .teamcity && mvn teamcity-configs:generate
