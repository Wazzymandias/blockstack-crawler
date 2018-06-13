#!/bin/sh

image="blockstack-watcher:latest"
docker run --name blockstack-watcher -d blockstack-watcher
