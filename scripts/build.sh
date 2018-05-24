#!/bin/bash

dep version || curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
dep ensure
cp data/config.go.sample data/config.go
go build main.go