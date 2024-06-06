#!/bin/bash
if [ "$#" -ne 1 ]
then
  echo "Usage: ./createFile.sh [TESTNAME_STRING]"
  exit 1
fi
if [ ! -d "/tmp/aztec" ]; then
  mkdir /tmp/aztec
fi
if [ ! -d "/tmp/aztec/$1" ]; then
  mkdir /tmp/aztec/$1
fi

cp ./circuits/astar/Nargo.toml /tmp/aztec/$1/Nargo.toml
cp ./circuits/astar/Verifier.toml /tmp/aztec/$1/Verifier.toml
cp -r ./circuits/astar/src /tmp/aztec/$1/
echo "" > /tmp/aztec/$1/Prover.toml && echo "File created"
