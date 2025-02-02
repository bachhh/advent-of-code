#!/bin/sh

set -eax
DAY=$1
cd $1
go build -o ./cmd/$1
# ./cmd/$1 -cpuprofile=profile.out
./cmd/$1
