#!/bin/bash

[ -z "${AWS_ACCESS_KEY_ID}" ] && echo "AWS_ACCESS_KEY_ID is required" && exit 1
[ -z "${AWS_SECRET_ACCESS_KEY}" ] && echo "AWS_SECRET_ACCESS_KEY is required" && exit 1

export res="$(/janitor $@)"

# https://github.community/t/set-output-truncates-multiline-strings/16852/2
res="${res//'%'/'%25'}"
res="${res//$'\n'/'%0A'}"
res="${res//$'\r'/'%0D'}"

echo "::set-output name=result::$(echo "$res")"