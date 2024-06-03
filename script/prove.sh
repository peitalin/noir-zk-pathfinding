#!/bin/bash
if [ "$#" -ne 1 ]
then
  echo "Usage: ./prove.sh [TESTNAME_STRING]"
  exit 1
fi

cd /tmp/aztec/$1 && nargo compile && nargo prove --oracle-resolver http://localhost:5555 && echo "Proof Generated"
