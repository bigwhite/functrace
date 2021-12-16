#!/bin/bash 

files=`find . -not -path "./vendor/**/*" -name "*.go" -print`

for f in $files
do
    echo $f
    gen -w $f
done
