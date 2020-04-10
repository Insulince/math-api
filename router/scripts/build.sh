#!/bin/bash

set -ex

: ========== BUILDING ROUTER ==========

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

releaseDockerfile="./Dockerfile" # The name of your Dockerfile for the release image.
buildContext="$mathAPiRootLocation" # What files Docker will consider when building (most likely the project root).
project="math-api" # Any sub path from the registry to your repository.
repository="router" # The name of your Docker repository.
releaseImage="$project/$repository:$tag" # The full image name for the release image.

# Pull the latest images to be used in the containers.
docker pull alexellis2/nginx-arm:latest

# Build the release-image.
docker build --no-cache -t "$releaseImage" -f "$releaseDockerfile" "$buildContext"

: Successfully built image "$releaseImage".
: ========== BUILD OF ROUTER SUCCESSFUL ==========
