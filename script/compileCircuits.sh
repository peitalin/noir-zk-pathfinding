#!/bin/bash
source .env

### compile circuits
nargo compile --workspace

### codegen solidity verifier contract
nargo codegen-verifier --workspace

### copy ./circuits/contract/astar/plonk_vk.sol to ./src/plonk_vk.sol
cp ./circuits/contract/astar/plonk_vk.sol ./src/plonk_vk.sol

### run prover and generate proof
nargo prove --oracle-resolver http://localhost:5555 --workspace
# TODO: need to find a way to feed inputs from Prover.toml into proving process programatically
# see `createFile.sh` and `prove.sh`

# TODO: need to feed public outputs from Verifier.toml into forge tests programmatically