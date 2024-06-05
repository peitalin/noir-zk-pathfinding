#!/bin/bash
source .env

forge create --rpc-url $RPC_URL \
    --private-key $PRIVATE_KEY \
    --etherscan-api-key $ETHERSCAN_APY_KEY \
    --verify \
    src/Astar.sol:Astar