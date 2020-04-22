#! /bin/bash

set -e

go build

./client  -addrs "localhost:8080,localhost:8081,localhost:8082"
