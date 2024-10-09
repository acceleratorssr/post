#!/bin/bash

outputDir="../output"

if [ ! -d "$outputDir" ]; then
    mkdir -p "$outputDir"
fi

services=("..\/article" "interactive" "search" "sso" "user" "bff")

for service in "${services[@]}"; do
    echo "Building $service..."
    cd "$service" || { echo "Failed to enter directory $service"; exit 1; }

    if ! go build -o "$outputDir" .; then
        echo "Error building $service"
        exit 1
    fi

    cd - || exit 1
done

cd "scripts" || exit 1
echo "All services built successfully!"
