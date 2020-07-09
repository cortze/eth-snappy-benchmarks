# eth-snappy-test

Benchmark code to meassure the compression ratio and speed of the snappy compressor(golang version).

The code reads all the json files of the given folder and start making the compression and decompression giving the average of the compress ratio and speed.

## Usage

1. git clone the repository
2. Build the go code
    ´go build -o snappy´
3. Run the code (by default reads the .json files in the /data folder)
    ´./snappy´
    
## Results

CPU used for the benchmark - Intel(R) Core(TM) i5-6600K CPU @ 3.50GHz

results compressing ethereum blocks:

´Average Compress Ratio:  1.8502550603952752

Max Ratio:  2.4160446800844184 Min Ratio:  1.5071121402580219

Average Compress Speed:  457.1678515328362 MB/s

Max Speed:  762.6247034738244 Min Speed:  15.078406170078727´
