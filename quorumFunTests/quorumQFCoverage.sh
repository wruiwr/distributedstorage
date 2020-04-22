#! /bin/bash

set -e


covDirName="../results/coverage"
resDirName="../results/qftests"

testDriverNames=("RD" "RDsqWR" "WRprRDsqRD" "WR" "WRprRD" "WRsqRD")
cov="ProfileCoverage"
resCov="CovTotal"

xmlQFDirName="../xml/qftest"

echo "running QFs coverage"


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
    go test -v -rdir="$xmlQFDirName/$tns/readqf.xml" -wdir="$xmlQFDirName/$tns/writeqf.xml" -coverprofile=$covDirName/$tns/qftest$cov.out > $covDirName/$tns/qftest$resCov.out

done

for tns in "${testDriverNames[@]}"
do
    cd $covDirName/$tns

    go tool cover -func=qftestProfileCoverage.out > qftestFuncCoverge.out

    cd ../../
done

echo "finished running QFs coverage"