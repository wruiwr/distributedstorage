#! /bin/bash

set -e

xmlDirName="../xml/qftest"
resDirName="../results/qftests"

testDriverNames=("RD" "RDsqWR" "WRprRDsqRD" "WR" "WRprRD" "WRsqRD" )


# create folders for results
if [ ! -d $resDirName ]; then
    for tns in "${testDriverNames[@]}"
    do
        mkdir -p $resDirName/$tns
    done
else
    echo "the directories for results exist"
fi

# tests
for tns in "${testDriverNames[@]}"
do
    go test -v qspecs_test.go qspecs.go -rdir="$xmlDirName/$tns/readqf.xml" -wdir="$xmlDirName/$tns/writeqf.xml"  > $resDirName/$tns/qfresult.out
done
