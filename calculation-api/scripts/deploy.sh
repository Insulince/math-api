#!/bin/bash

set -ex

: ========== DEPLOYING CALCULATION-API ==========

if [ "$1" == "" ]; then
    : ERROR: Missing first command line argument: Release version.
    exit "1"
fi
version="$1" # The location of the math-api root directory.

helmChartLocation="$HELM_CHARTS_ROOT/math-api/calculation-api"
chartLocation="$helmChartLocation/Chart.yaml"
valuesLocation="$helmChartLocation/values.yaml"

sed -E "s/version: .+/version: $version/" "$chartLocation" > "$chartLocation.intermediate"
mv "$chartLocation.intermediate" "$chartLocation"

sed -E "s/  tag: .+/  tag: $version/" "$valuesLocation" > "$valuesLocation.intermediate"
mv "$valuesLocation.intermediate" "$valuesLocation"

helm upgrade --install "calculation-api" "$helmChartLocation"

: Successfully deployed release "calculation-api".
: ========== DEPLOYMENT OF CALCULATION-API SUCCESSFUL ==========
