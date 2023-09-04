#!/bin/bash

ETH_NODE="localhost:5052"
ENCODING="ssz"
BLOCK_LIST="target_blocks.csv"
OUTPUT_FOLDER=$ENCODING"_blocks"

./build/snappy-benchmark fetch-blocks \
    --eth-node $ETH_NODE  \
    --encoding $ENCODING \
    --block-list $BLOCK_LIST \
    --output-folder $OUTPUT_FOLDER
