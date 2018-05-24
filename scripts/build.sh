#!/bin/bash

printenv
dep version || curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
dep ensure
go build main.go -o main