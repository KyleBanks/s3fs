#!/bin/bash
#
# Runs `golint` on all non-vendor packages.
#
# Usage: 
# 	./lint.sh

packages=$(go list ./... | grep -v vendor)

for p in $packages; do

	out=$(golint $p)
	if [[ !  -z  $out  ]]; then 
		echo $p
		echo "-------------------------------------------"
		echo ""
		echo $out
		echo ""	
	fi
	
done	