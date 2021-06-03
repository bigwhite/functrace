#!/bin/bash 

files=`find . -name "*.go" -print`
echo $files

for f in $files
do 
		gen -w $f
done
