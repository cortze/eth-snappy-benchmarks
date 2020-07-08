package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang/snappy"
)

func main() {
	var dataPath = flag.String("path", "/data", "Path of the data for the snappy compression test")
	var suffix = flag.String("suffix", ".json", "Path of the data for the snappy compression test")
	flag.Parse()

	var files []string
	var avgCompresRatio []float64
	var avgCompresSpeed []float64

	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	// Add to working directory the path of the folder where all the test are
	wd = wd + *dataPath

	// Get all the names of the json files on the data folder
	err = filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Iterate through the json files to make the compression test
	for _, file := range files {

		fmt.Println(file)

		// filter the files that finishes in .json
		if strings.HasSuffix(file, *suffix) {

			// Open our jsonFile
			jsonFile, err := os.Open(file)

			// if we os.Open returns an error then handle it
			if err != nil {
				fmt.Println(err)
			} else {

				// defer the closing of our jsonFile so that we can parse it later on
				// Read the opened jsonFile as a byte array
				byteValue, _ := ioutil.ReadAll(jsonFile)
				fmt.Println("Size block (bytes):            ", len(byteValue))

				// Start the timer
				start1 := time.Now()

				// Compress the message with the snappy compressor
				compressmsg := snappy.Encode(nil, byteValue)
				if err != nil {
					fmt.Println("Encode Failed")
				}

				fmt.Println("Size block compressed (bytes): ", len(compressmsg))

				// Stop the timer
				codetime := time.Since(start1)
				ctime := float64(codetime*time.Microsecond) / float64(time.Millisecond)

				fmt.Println("Encoding time:     ", codetime)

				// Run the timer
				start2 := time.Now()

				//fmt.Println("Message Compressed: ", compressmsg.Data)

				// Decode the message
				_, err = snappy.Decode(nil, compressmsg)
				if err != nil {
					fmt.Println("Decode Failed")
				}

				decodetime := time.Since(start2)

				compressRatio := (float64(len(byteValue)) / float64(len(compressmsg)))
				compressSpeed := (float64(len(byteValue)) / ctime)

				fmt.Println("Decoding time:     ", decodetime)
				fmt.Println("Compression ratio: ", compressRatio)
				fmt.Println("Compression Speed: ", compressSpeed, "MB/s")

				avgCompresRatio = append(avgCompresRatio, compressRatio)
				avgCompresSpeed = append(avgCompresSpeed, compressSpeed)

				fmt.Printf("\n")
			}

			defer jsonFile.Close()
		}
	}
	var avgRatio float64 = 0
	var avgSpeed float64 = 0

	// Run stadistics of the taken averages
	for i := 0; i < len(avgCompresRatio); i++ {
		avgRatio = avgRatio + avgCompresRatio[i]
		avgSpeed = avgSpeed + avgCompresSpeed[i]
	}
	fmt.Println("Average Compress Ratio: ", avgRatio/float64(len(avgCompresRatio)))
	fmt.Println("Average Compress Speed: ", avgSpeed/float64(len(avgCompresSpeed)), "MB/s")

}
