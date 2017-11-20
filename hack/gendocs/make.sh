#!/usr/bin/env bash

pushd $GOPATH/src/github.com/k8sdb/xdb/hack/gendocs
go run main.go
popd
