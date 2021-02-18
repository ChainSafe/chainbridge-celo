#!/usr/bin/env bash

CVPKG=$(go list ./... | grep -v 'e2e\|generated\|bindata\|mock\|main.go\|bindings\|shared\|root.go' | tr '\n' ',')
go test -coverpkg=$CVPKG -coverprofile=cover.out -p=1 $(go list ./... | grep -v 'e2e')
