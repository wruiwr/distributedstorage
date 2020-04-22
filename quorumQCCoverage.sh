#! /bin/bash

set -e

#addr="localhost:8080,localhost:8081,localhost:8082"
xmlDirName="./xml/systest"
resDirName="./results/systests"
covDirName="./results/coverage"

testDriverNames=("RD" "RDsqWR" "WRprRDsqRD" "WR" "WRprRD" "WRsqRD" "ALL")
cov="coverage"
resCov="covTotal"

# create folders for results
if [ ! -d $resDirName ]; then
    for tns in "${testDriverNames[@]}"
    do
        mkdir -p $resDirName/$tns
        mkdir -p $covDirName/$tns
    done
elif [ ! -d $covDirName ]; then
    for tns in "${testDriverNames[@]}"
    do
        mkdir -p $covDirName/$tns
    done
else
    echo "the directories for results exist"
fi

# coverage QCs
for tns in "${testDriverNames[@]}"
do
    for filename in $xmlDirName/$tns/systemtest*
    do
        cns=$(basename ${filename%.*})
        echo $tns/$cns
        #sh ./run3servers.sh
        go test -dir="$filename" -coverprofile=$covDirName/$tns/$cov$cns.out > $covDirName/$tns/$resCov$cns.out
        #sh ./stop3servers.sh
    done
done


for tns in "${testDriverNames[@]}"
do
    cd $covDirName/$tns
    for filename in $cov*
    do
        go tool cover -func=$filename > covQctest$filename
    done
    cd ../../..
done








