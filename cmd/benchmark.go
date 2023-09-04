package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cortze/eth-snappy-benchmarks/metrics"

	"github.com/golang/snappy"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	blocksFolder     string = ""
	metricsFolder    string = ""
	blocksEncoding   string = ""
	blocksIterations int    = 10
)

var RunCmd = &cli.Command{
	Name:        "run",
	Description: "",
	Action:      RunBenchmark,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "block-folder",
			DefaultText: "folder placing all the blocks to test",
			Required:    true,
			Destination: &blocksFolder,
		},
		&cli.StringFlag{
			Name:        "metrics-folder",
			DefaultText: "folder placing all the final csvs from the benchmarks",
			Required:    true,
			Destination: &metricsFolder,
		},
		&cli.StringFlag{
			Name:        "encoding",
			DefaultText: "folder placing all the blocks to test",
			Required:    true,
			Destination: &blocksEncoding,
		},
		&cli.IntFlag{
			Name:        "iterations",
			Aliases:     []string{"i"},
			DefaultText: "folder placing all the blocks to test",
			Destination: &blocksIterations,
		},
	},
}

func RunBenchmark(ctx *cli.Context) error {
	log.WithFields(log.Fields{
		"block-folder":     blocksFolder,
		"metrics-folder":   metricsFolder,
		"block-suffix":     blocksEncoding,
		"block-iterations": blocksIterations,
	}).Info("launching snappy compession benchmarks")

	files, err := getFilesFromFolder(blocksFolder, blocksEncoding)
	if err != nil {
		return errors.Wrap(err, "unable to get errors from folder "+blocksFolder)
	}

	log.WithFields(log.Fields{
		"folder":                    blocksFolder,
		(blocksEncoding + "-files"): len(files),
	}).Info("found...")

	// Iterate through the json files to make the compression test
	for _, file := range files {

		cleanFile := strings.TrimPrefix(file, blocksFolder+"/")

		metricsFile := strings.Split(cleanFile, ".")[0] + ".csv"
		compressionMetrics, err := metrics.NewCompressionMetrics(blocksFolder, metricsFolder, metricsFile, blocksEncoding)
		if err != nil {
			log.WithFields(log.Fields{
				"metrics-file": metricsFile,
			}).Error("unable to create metrics file")
			continue
		}

		fileBytes, err := fileToBytes(file)
		if err != nil {
			log.WithFields(log.Fields{
				"file":  cleanFile,
				"error": err,
			}).Error("unable to open file")
			continue
		}

		// Run the compression test 10 times for each block
		for i := 0; i < blocksIterations; i++ {

			// --- Compression Starts ---
			startTime := time.Now()
			compressMsg := snappy.Encode(nil, fileBytes)
			encodeTime := time.Since(startTime)

			// --- Decompression Starts ---
			startTime = time.Now()
			_, err = snappy.Decode(nil, compressMsg)
			decodeTime := time.Since(startTime)
			if err != nil {
				fmt.Println("Decode Failed")
				continue
			}

			// compression metrics
			compressRatio := float64(len(fileBytes)) / float64(len(compressMsg))
			compressSpeed := float64(len(fileBytes)) / float64(encodeTime.Nanoseconds())

			// basic info on screen
			if i == 0 {
				log.Infof("%s - Block size(Bytes) %d - Compressed Block size (bytes) %d\n", cleanFile, len(fileBytes), len(compressMsg))
				log.Info("Encoding time; Decoding time; Compression ratio ; Compression Speed (Bytes/Millisecond)\n")
			}
			log.WithFields(log.Fields{
				"encoding-time":  encodeTime,
				"decoding-time":  decodeTime,
				"compress-ratio": compressRatio,
				"compress-speed": compressSpeed,
			}).Info("")

			compressionMetrics.AddResults(
				cleanFile,
				int64(len(fileBytes)),
				int64(len(compressMsg)),
				encodeTime,
				decodeTime,
				compressRatio,
				compressSpeed)
		}

		summary := compressionMetrics.GetSummary()
		log.WithFields(log.Fields{
			"encoding-time":  summary[metrics.Minimum][metrics.EncodingTime],
			"decoding-time":  summary[metrics.Minimum][metrics.DecodingTime],
			"compress-ratio": summary[metrics.Minimum][metrics.CompressRatio],
			"compress-speed": summary[metrics.Minimum][metrics.CompressSpeed],
		}).Info("MIN")

		log.WithFields(log.Fields{
			"encoding-time":  summary[metrics.Maximum][metrics.EncodingTime],
			"decoding-time":  summary[metrics.Maximum][metrics.DecodingTime],
			"compress-ratio": summary[metrics.Maximum][metrics.CompressRatio],
			"compress-speed": summary[metrics.Maximum][metrics.CompressSpeed],
		}).Info("MAX")

		log.WithFields(log.Fields{
			"encoding-time":  summary[metrics.Average][metrics.EncodingTime],
			"decoding-time":  summary[metrics.Average][metrics.DecodingTime],
			"compress-ratio": summary[metrics.Average][metrics.CompressRatio],
			"compress-speed": summary[metrics.Average][metrics.CompressSpeed],
		}).Info("AVG")

		err = compressionMetrics.Export()
		if err != nil {
			log.Error(errors.Wrap(err, "unable to export metrics to csv file"))
		}
	}
	return nil
}

func fileToBytes(file string) ([]byte, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return []byte{}, err
	}
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		return []byte{}, err
	}
	return fileBytes, nil
}

func getFilesFromFolder(folder, suffix string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(folder, func(file string, info os.FileInfo, err error) error {
		if strings.HasSuffix(file, suffix) {
			files = append(files, file)
		}
		return nil
	})
	return files, err
}
