#!/bin/bash

files=`find . -not -path "./vendor/**/*" -name "*.go" -print`

echo $files

for f in $files
do
    gen -w $f
done
