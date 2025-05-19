#!/usr/bin/env bash
set -eoux pipefail


curl -X POST \
    -H "Cookie: session=${COOKIE}" \
    -d "level=${LEVEL}&answer=${ANSWER}" \
    https://adventofcode.com//day/${DAY}/answer

# year is the last part of PWD
YEAR=$(basename $PWD)
# parse day, level from input
IFS="-" read DAY LEVEL <<< "$1"
ANSWER=$2

# Endpoint URL
URL="https://adventofcode.com/${YEAR}/day/${DAY}/answer"

# Submit the answer via curl
curl -s -X POST "$URL" \
    -H "Cookie: session=${COOKIE}" \
    -d "level=${PART}" \
    -d "answer=${ANSWER}" \
    -H "Content-Type: application/x-www-form-urlencoded"
