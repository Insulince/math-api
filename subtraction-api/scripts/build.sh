#!/bin/bash

set -ex

: ========== BUILDING SUBTRACTION-API ==========

if [ "$1" == "" ]; then
    : ERROR: Missing first command line argument: Image Tag.
    exit "1"
fi
tag="$1" # The tag for this deployment.

if [ "$2" == "" ]; then
    : ERROR: Missing second command line argument: Math API root location.
    exit "1"
fi
mathAPiRootLocation="$2" # The location of the math-api root directory.

buildDockerfile="./Dockerfile.build" # The name of your Dockerfile for the build image.
releaseDockerfile="./Dockerfile.release" # The name of your Dockerfile for the release image.
buildContext="$mathAPiRootLocation" # What files Docker will consider when building (most likely the project root).
project="math-api" # Any sub path from the registry to your repository.
repository="subtraction-api" # The name of your Docker repository.
buildImage="$project/$repository:build" # The full image name for the build image.
releaseImage="$project/$repository:$tag" # The full image name for the release image.
buildArtifact="app" # The name of the artifact generated in the build image for use in the release image.
buildArtifactLocation="/go/src/$project/$repository/bin/$buildArtifact" # The location of the build-artifact in the build-image.

# Pull the latest images to be used in the containers.
docker pull golang:alpine
docker pull alpine:latest

# Build the build-image.
docker build -t "$buildImage" -f "$buildDockerfile" "$buildContext"

# Extract the build-artifact from the build-image.
docker create --name build_container "$buildImage"
docker cp build_container:"$buildArtifactLocation" "./$buildArtifact"
docker rm -f build_container

# Build the release-image.
docker build -t "$releaseImage" -f "$releaseDockerfile" "$buildContext"

# Clean up build process leftovers.
rm -r "./$buildArtifact"
docker rmi  "$buildImage"

: Successfully built image "$releaseImage".
: ========== BUILD OF SUBTRACTION-API SUCCESSFUL ==========
