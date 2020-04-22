#! /bin/bash

set -e

go build

./server -port=8080 &

./server -port=8081 &

./server -port=8082 &

echo "running, enter to stop"

read && killall server
