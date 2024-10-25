#!/bin/bash

# Directories to exclude from testing
EXCLUDE_DIRS="(proto|vendor|doc|mock|tool|scripts|model|gen|test|migrations)"

# Generate a list of all go packages, exclude the specified directories
PACKAGES=$(go list ./... | grep -vE "$EXCLUDE_DIRS")

# Run go test on the selected packages
go test -coverprofile=doc/test-report-coverage.out $PACKAGES

# Generate a coverage report
go tool cover -html=doc/test-report-coverage.out -o doc/test-report-coverage.html
