#!/bin/bash
# Coverage test script for config package

set -e

echo "Running tests with coverage reporting..."
go test -v -coverprofile=coverage.out -covermode=atomic ./internal/config

echo ""
echo "Converting coverage to HTML format..."
go tool cover -html=coverage.out -o coverage.html

echo ""
echo "Generating coverage summary..."
go tool cover -func=coverage.out | tail -20

echo ""
echo "Coverage report generated:"
echo "  - coverage.out (raw coverage data)"
echo "  - coverage.html (HTML report)"
