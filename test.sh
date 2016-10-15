#!/bin/bash
#
# Runs `go test` on all non-vendor packages.
#
# Usage: 
# 	./test.sh

go test -cover $(go list ./... | grep -v vendor)