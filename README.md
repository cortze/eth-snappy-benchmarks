# eth-snappy-test

Benchmark code to meassure the compression ratio and speed of the snappy compressor(golang version).

The code reads all the json files of the given folder and start making the compression and decompression giving the average of the compress ratio and speed.

## Usage

1. git clone the repository
2. Build the go code
    ´go build -o snappy´
3. Run the code (by default reads the .json files in the /data folder)
    ´./snappy´
