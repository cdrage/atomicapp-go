#!/bin/bash

# For now we only build on linux
go get $DEP_ARGS gopkg.in/yaml.v1
go get $DEP_ARGS github.com/Sirupsen/logrus
go get $DEP_ARGS github.com/codegangsta/cli

# Get the rest of the deps
DEPS=$(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
go get $DEP_ARGS ./... $DEPS
