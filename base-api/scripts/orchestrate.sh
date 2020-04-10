#!/bin/bash

set -ex

: ============================== ORCHESTRATING MATH-API ==============================

if [ "$1" == "" ]; then
    : ERROR: Missing first command line argument: Docker Registry.
    exit "1"
fi
registry="$1" # Location of Docker registry.

if [ "$2" == "" ]; then
    : ERROR: Missing second command line argument: Docker Image version.
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
mathApiRootLocation="$4" # The location of the math-api root directory.

originalLocation="$(pwd)"
projectNames=(
    "router"
    "calculation-api"
    "addition-api"
    "subtraction-api"
    "multiplication-api"
    "division-api"
)

for projectName in "${projectNames[@]}"; do
    cd "$mathApiRootLocation/$projectName"
    bash "./scripts/orchestrate.sh" "$registry" "$tag" "$orchestrationLevel" "$mathApiRootLocation" &
    cd "$originalLocation"
done

: ============================== ORCHESTRATION OF MATH-API SUCCESSFUL ==============================
