#!/bin/bash
source .env

### compile circuits
nargo compile

### codegen solidity verifier contract
nargo codegen-verifier

### copy ./circuits-radius/contract/plonk_vk.sol to ./src/plonk_vk.sol
cp ./circuits-radius/contract/plonk_vk.sol ./src/plonk_vk.sol


### run prover and generate proof
nargo prove
# TODO: need to find a way to feed inputs from Prover.toml into proving process programatically

### copy public inputs into tests
# TODO: need to feed public outputs from Verifier.toml into forge tests programmatically