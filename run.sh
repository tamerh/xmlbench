#!/bin/sh

if [[ ! -f uniprot_sprot.xml.gz ]]; then
    curl -OL https://ftp.uniprot.org/pub/databases/uniprot/current_release/knowledgebase/complete/uniprot_sprot.xml.gz
fi


rm -f xmlbench

go build xmlbench.go 

echo "starting prod"
# to profile run with profile param
# ./xmlbench -profile
./xmlbench
echo "killing top process"
pkill -f "top -stats"

echo "starting development version"
# ./xmlbench -devel -profile
./xmlbench -devel
echo "killing top process"
pkill -f "top -stats"

# clean manually the top results. e.g using following regex to remove lines 
# Processes.*\n
# 2020/12/17.*\n
# Load.*\n
# CPU usage.*\n
# SharedLibs.*\n
# MemRegions.*\n
# PhysMem.*\n
# VM:.*\n
# Networks:.*\n
# Disks:.*\n\n
# %CPU.*\n



# to view profiler output use pprof -http=localhost:8888
