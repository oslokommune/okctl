#!/bin/sh -l

res=$(./janitor $1)
echo "::set-output name=result::$res"