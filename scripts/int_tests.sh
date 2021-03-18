#!/usr/bin/env bash
go test -p=1 $(go list ./... | grep 'cbcli\|e2e')