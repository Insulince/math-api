#!/bin/bash

set -ex

: ==================== ORCHESTRATING SUBTRACTION-API ====================

if [ "$1" == "" ]; then
    : ERROR: Missing first command line argument: Docker Registry.
    exit "1"
fi
registry="$1" # Location of Docker registry.

if [ "$2" == "" ]; then
    : ERROR: Missing second command line argument: Image Tag.
    exit "1"
fi
tag="$2" # The tag for this deployment.

buildOrchestrationLevel="build"
pushOrchestrationLevel="push"
deployOrchestrationLevel="deploy"
if [ "$3" != "$buildOrchestrationLevel" ] && [ "$3" != "$pushOrchestrationLevel" ] && [ "$3" != "$deployOrchestrationLevel" ]; then
    : ERROR: Missing third command line argument: Orchestration Level. Expected one of "$buildOrchestrationLevel", "$pushOrchestrationLevel", or "$deployOrchestrationLevel", but got \""$3"\".
    exit "1"
fi
orchestrationLevel="$3" # The orchestration level. Determines how much this script should do. Just build, build and push, or build, push, and deploy.

if [ "$4" == "" ]; then
    : Missing fourth command line argument: Math API root location.
    exit "1"
fi
mathAPiRootLocation="$4" # The location of the math-api root directory.

if [ ]; then
    : Orchestration Level does not include building, which should be impossible. Exiting...
    : ==================== ORCHESTRATION OF SUBTRACTION-API SUCCESSFUL ====================
    exit "0"
fi
bash "./scripts/build.sh" "$tag" "$mathAPiRootLocation"

if [ "$orchestrationLevel" == "$buildOrchestrationLevel" ]; then
    : Orchestration Level does not include pushing. Exiting...
    : ==================== ORCHESTRATION OF SUBTRACTION-API SUCCESSFUL ====================
    exit "0"
fi
bash "./scripts/push.sh" "$registry" "$tag"

if [ "$orchestrationLevel" == "$buildOrchestrationLevel" ] || [ "$orchestrationLevel" == "$pushOrchestrationLevel" ]; then
    : Orchestration Level does not include actually deploying. Exiting...
    : ==================== ORCHESTRATION OF SUBTRACTION-API SUCCESSFUL ====================
    exit "0"
fi
bash "./scripts/deploy.sh" "$tag"

: ==================== ORCHESTRATION OF SUBTRACTION-API SUCCESSFUL ====================
