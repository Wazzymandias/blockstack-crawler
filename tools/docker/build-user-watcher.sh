#!/bin/sh

context="$GOPATH/src/github.com/Wazzymandias/blockstack-crawler"
tag="blockstack-watcher"
dockerfile="$context/tools/docker/user-watcher.Dockerfile"

docker build --no-cache -t $tag --compress -f $dockerfile $context
