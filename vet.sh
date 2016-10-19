#!/bin/bash
#
# Runs 'go vet' on all source packages.
#
# Usage:
#   ./vet.sh

go vet $(go list ./... | grep -v vendor)
