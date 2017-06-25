#!/bin/bash

mkdir -p $1/config
cd $1 && go build github.com/awgh/marvin/marvin && cp $GOPATH/src/github.com/awgh/marvin/db/* . && cp $GOPATH/src/github.com/awgh/marvin/config/* .

