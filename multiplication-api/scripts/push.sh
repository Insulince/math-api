#!/bin/bash

set -ex

: ========== PUSHING MULTIPLICATION-API ==========

if [ "$1" == "" ]; then
    : ERROR: Missing first command line argument: Docker Registry.
    exit "1"
fi
registry="$1"

if [ "$2" == "" ]; then
    : ERROR: Missing second command line argument: Image Tag.
    exit "1"
fi
tag="$2"

image="math-api/multiplication-api:$tag"

taggedImage="$registry/$image"

docker tag "$image" "$taggedImage"
docker push "$taggedImage"

: Successfully pushed image "$taggedImage".
: ========== PUSH OF MULTIPLICATION-API SUCCESSFUL ==========
