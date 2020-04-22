# README

This README would briefly explain the usage of the Gorums rwregister.

## To regenerate the `.pb.go` file

```sh
go generate
```

## To run the tests

```sh
go test -v
```

## To run server and client only

1. startserver.sh in server folder can start 3 servers.

1. startclient.sh in client folder can run a client.

1. Need to start servers first

## Run whole system and tests

* quorumcalltests.sh in registerQspecs can run tests for the system.
* quorumQCCoverage.sh in registerQspecs can run tests for the system and get coverages.
* stop3servers.sh can stop servers.
* results will be stored in results folder.

## Run quorum functions tests

* quorumQFtests.sh can run the tests for quorum functions
* quorumQFCoverage.sh can run the tests and get coverages of quorum functions.

## Other information

xmlReader.go in reader is the Go code to read test cases from xml files.

Quorum functions (ReadQF and WriteQF) are implemented in qspecs.go in the registerQspecs.

quorumcalltests.sh and quorumQFtests.sh can start tests.
coverage.sh can run an example of the test and get the result of coverage.
