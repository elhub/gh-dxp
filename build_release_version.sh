#!/bin/bash

# Set the plugin name
PLUGIN_NAME="gh-dxp"

# Create the dist directory if it doesn't exist
mkdir -p dist

# Platforms to build for
PLATFORMS=("linux/amd64" "darwin/arm64")

# Function to compile the plugin for a specific platform
compile() {
  local GOOS=$1
  local GOARCH=$2
  local OUTPUT="dist/${PLUGIN_NAME}-${GOOS}-${GOARCH}"

  echo "Building for ${GOOS}/${GOARCH}..."

  GOOS=${GOOS} GOARCH=${GOARCH} go build -o "${OUTPUT}" .

  EXIT_CODE=$?
  
  if [ $EXIT_CODE -ne 0 ]; then
    echo "Error building for ${GOOS}/${GOARCH}"
    exit 1
  fi
}

# Iterate over platforms and compile for each
for PLATFORM in "${PLATFORMS[@]}"; do
  IFS="/" read -r GOOS GOARCH <<< "${PLATFORM}"
  compile "${GOOS}" "${GOARCH}"
done

echo "All builds completed successfully."
