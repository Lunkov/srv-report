#!/bin/bash

go get -u github.com/rakyll/gotest
gotest -v -covermode=count -coverprofile=coverage.out .
go tool cover -func=coverage.out
go tool cover -html=coverage.out
