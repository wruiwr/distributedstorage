#! /bin/bash

set -e

#addr="localhost:8080,localhost:8081,localhost:8082"
xmlDirName="./xml/systest"
resDirName="./results/systests"
#covDirName="../results/coverage"

testDriverNames=("RD" "RDsqWR" "WRprRDsqRD" "WR" "WRprRD" "WRsqRD" "ALL")
#cov="coverage"
#resCov="covresult"

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
    for filename in $xmlDirName/$tns/systemtest*
    do
        cns=$(basename ${filename%.*})
        echo $tns/$cns
        #sh ./run3servers.sh
        go test -v rwregister_test.go rwregister.go rwqspec.go register.pb.go monitor.go -dir="$filename" > $resDirName/$tns/$cns.out
        #sh ./stop3servers.sh
    done
done

