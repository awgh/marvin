#!/bin/bash

GP=`go env GOPATH`

mkdir -p $1/config
go get github.com/mattn/go-sqlite3 && go get github.com/fluffle/goirc && go get github.com/awgh/markov && go get github.com/slack-go/slack && go get github.com/gocolly/colly/...
cd $1 && go build github.com/awgh/marvin && cp $GP/src/github.com/awgh/marvin/db/* . && cp $GP/src/github.com/awgh/marvin/config/* .

