#!/bin/sh

rm -f $2.txt
top -stats cpu,idlew,power -o power -d -pid $1 > $2.txt