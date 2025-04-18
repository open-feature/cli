#!/bin/bash
set -e

# Script to test if the generated C# code compiles correctly
SCRIPT_DIR=$(dirname "$0")
CLI_ROOT=$(realpath "$SCRIPT_DIR/../..")
OUTPUT_DIR=$(realpath "$SCRIPT_DIR")

echo "=== Building OpenFeature CLI ==="
cd "$CLI_ROOT"
go build

echo "=== Generating C# client ==="
./cli generate csharp --manifest="$CLI_ROOT/sample/sample_manifest.json" --output="$OUTPUT_DIR/expected" --namespace="TestNamespace"

if [ ! -f "$OUTPUT_DIR/expected/OpenFeature.cs" ]; then
    echo "Error: OpenFeature.cs was not generated"
    exit 1
fi

echo "=== Building Docker image to compile C# code ==="
cd "$OUTPUT_DIR"
docker build -t openfeature-csharp-test .

echo "=== Testing C# compilation and execution ==="
docker run --rm openfeature-csharp-test

if [ $? -eq 0 ]; then
    echo "=== Success: C# code compiles and executes correctly ==="
    exit 0
else
    echo "=== Error: C# code fails to compile or execute ==="
    exit 1
fi