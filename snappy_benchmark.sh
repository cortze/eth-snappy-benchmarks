#!/bin/bash

ENCODING="ssz"
BLOCK_FOLDER=$ENCODING"_blocks"
METRICS_FOLDER="results"
ITERATIONS=10

./build/snappy-benchmark run \
  --block-folder $BLOCK_FOLDER \
  --metrics-folder $METRICS_FOLDER \
  --encoding $ENCODING \
  --iterations $ITERATIONS
